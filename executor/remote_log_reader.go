// Copyright 2018 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package executor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/pingcap/errors"
	"github.com/pingcap/parser/model"
	"github.com/pingcap/parser/mysql"
	log2 "github.com/pingcap/tidb/infoschema/inspection/log"
	"github.com/pingcap/tidb/types"
	"github.com/pingcap/tidb/util/chunk"
	"github.com/sirupsen/logrus"

	"golang.org/x/net/context"
)

// make sure `LocalLogReaderExecutor` implements `Executor`.
var _ Executor = &RemoteLogReaderExecutor{}

type RemoteLogReaderExecutor struct {
	baseExecutor
	Table   *model.TableInfo
	Columns []*model.ColumnInfo

	startTimeStr string
	endTimeStr string
	LimitStr string
	result *chunk.Chunk
	cnt int
	limit int
	pattern string
	url string
	fd string
}

func (e *RemoteLogReaderExecutor) Open(ctx context.Context) error {

	TimeStampLayout := "2006-01-02T15:04:05"
	local, err := time.LoadLocation("Asia/Chongqing")
	if err != nil {
		return err
	}
	// startTime, _ := time.ParseInLocation(TimeStampLayout ,"1970-01-01T00:00:00", local)
	// endTime, _ := time.ParseInLocation(TimeStampLayout , "2030-01-01T00:00:00", local)
	_, err = time.ParseInLocation(TimeStampLayout ,e.startTimeStr, local)
	if err != nil {
		return err
	}
	_, err = time.ParseInLocation(TimeStampLayout , e.endTimeStr, local)
	if err != nil {
		return err
	}
	l, err := strconv.ParseInt(e.LimitStr, 10, 64)
	if err != nil {
		return err
	}
	e.limit = int(l)
	e.cnt = 0
	url := fmt.Sprintf("%s/open?startTime=%s&endTime=%s", e.url, e.startTimeStr, e.endTimeStr)
	logrus.Infof("OpenURL=%s",url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	v := make(map[string]string)
	if err := json.Unmarshal(body, &v); err != nil {
		return err
	}
	var exists bool
	e.fd, exists = v["fd"]
	if !exists {
		return errors.New("illegal http return. missing field: `fd`")
	}
	logrus.Infof("FD=%s",e.fd)
	return nil
}

// Next fills data into the chunk passed by its caller.
func (e *RemoteLogReaderExecutor) Next(ctx context.Context, chk *chunk.Chunk) error {
	var err error

	chk.GrowAndReset(chk.Capacity())
	if e.result == nil {
		e.result = newFirstChunk(e)
	}
	e.result.Reset()
	err = e.fetchAll()
	if err != nil {
		return errors.Trace(err)
	}
	iter := chunk.NewIterator4Chunk(e.result)
	for colIdx := 0; colIdx < e.Schema().Len(); colIdx++ {
		retType := e.Schema().Columns[colIdx].RetType
		if !types.IsTypeVarchar(retType.Tp) {
			continue
		}
		for row := iter.Begin(); row != iter.End(); row = iter.Next() {
			if valLen := len(row.GetString(colIdx)); retType.Flen < valLen {
				retType.Flen = valLen
			}
		}
	}
	chk.Append(e.result, 0, e.result.NumRows())
	return nil
}

// Close implements the Executor Close interface.
func (e *RemoteLogReaderExecutor) Close() error {
	/*
	url := fmt.Sprintf("%s/close?fd=%s", e.url, e.fd)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	 */
	return nil
}

func (e *RemoteLogReaderExecutor) fetchAll() error {
	err := e.fetchRemoteLog()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (e *RemoteLogReaderExecutor) fetchRemoteLog() error {
	url := fmt.Sprintf("%s/next?fd=%s&limit=%d", e.url, e.fd, e.limit)
	logrus.Infof("URL=%s|",url)
	if e.pattern != "" {
		url = fmt.Sprintf("%s&pattern=%s", url, e.pattern)
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	logrus.Infof("FD=%s", e.fd)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Info("ERR=%s", err)
		return err
	}
	logrus.Info("RET=%s", string(body))
	logBatch := &log2.TiDBLogBatch{}
	if err := json.Unmarshal(body, logBatch); err != nil {
		return err
	}
	for _, item := range logBatch.Logs {
		if e.cnt >= e.limit {
			break
		}
		if item == nil {
			continue
		}
		data, err := e.getData(item)
		if err != nil {
			return errors.Trace(err)
		}
		row := chunk.MutRowFromDatums(data).ToRow()
		e.result.AppendRow(row)
		e.cnt++
	}
	return nil
}

func (e *RemoteLogReaderExecutor) getData(data *log2.TiDBLogItem) ([]types.Datum, error) {
	row := make([]types.Datum, 0, len(e.Columns))
	for _, col := range e.Columns {
		switch col.Name.L {
		case "host":
			row = append(row, types.NewStringDatum(data.Host))
		case "port":
			row = append(row, types.NewStringDatum(data.Port))
		case "component":
			row = append(row, types.NewStringDatum(data.Component))
		case "filename":
			row = append(row, types.NewStringDatum(data.FileName))
		case "time":
			tm := types.Time{
				Time: types.FromGoTime(data.Time),
				Type: mysql.TypeDatetime,
				Fsp:  0,
			}
			row = append(row, types.NewTimeDatum(tm))
		case "level":
			row = append(row, types.NewIntDatum(int64(data.Level)))
		case "content":
			row = append(row, types.NewStringDatum(data.Content))
		default:
			data := types.NewDatum(nil)
			row = append(row, data)
		}
	}
	return row, nil
}
