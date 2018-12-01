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
	"github.com/cznic/mathutil"
	"github.com/pingcap/errors"
	"github.com/pingcap/parser/model"
	"github.com/pingcap/tidb/types"
	"github.com/pingcap/tidb/types/json"
	"github.com/pingcap/tidb/util/chunk"
	"github.com/pingcap/tidb/util/mock"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// make sure `StreamReaderExecutor` implements `Executor`.
var _ Executor = &StreamReaderExecutor{}

var batchFetchCnt = 10
var maxFetchCnt = 10000
var streamCursor = 0

// streamTableMap : {type|topic|table =>pos}
var streamTableMaps = make(map[string]int64)

// StreamReaderExecutor reads data from a stream.
type StreamReaderExecutor struct {
	baseExecutor
	Table   *model.TableInfo
	Columns []*model.ColumnInfo

	isInit bool
	result *chunk.Chunk
	cursor int
}

func (e *StreamReaderExecutor) init() error {
	if e.isInit {
		return nil
	}

	e.isInit = true
	return nil
}

// Open initialzes necessary variables for using this executor.
func (e *StreamReaderExecutor) Open(ctx context.Context) error {
	err := e.init()
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

// Next fills data into the chunk passed by its caller.
func (e *StreamReaderExecutor) Next(ctx context.Context, chk *chunk.Chunk) error {
	log.Warnf("[qiuyesuifeng]%v:%v", e.Table, e.cursor)

	chk.GrowAndReset(e.maxChunkSize)
	if e.result == nil {
		e.result = e.newFirstChunk()
		err := e.fetchAll()
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
	}
	if e.cursor >= e.result.NumRows() {
		return nil
	}
	numCurBatch := mathutil.Min(chk.Capacity(), e.result.NumRows()-e.cursor)
	chk.Append(e.result, e.cursor, e.cursor+numCurBatch)
	e.cursor += numCurBatch
	streamCursor += numCurBatch
	return nil
}

// Close implements the Executor Close interface.
func (e *StreamReaderExecutor) Close() error {
	e.isInit = false
	return nil
}

func (e *StreamReaderExecutor) fetchAll() error {
	err := e.fetchMockData()
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (e *StreamReaderExecutor) fetchMockData() error {
	log.Warnf("[qiuyesuifeng][fetch mock data]%v", e.cursor)

	for i := streamCursor; i < maxFetchCnt && i < streamCursor+batchFetchCnt; i++ {
		row := []interface{}{mock.MockStreamData[i].ID, mock.MockStreamData[i].Content, mock.MockStreamData[i].CreateTime}
		e.appendRow(e.result, row)
	}

	return nil
}

func (e *StreamReaderExecutor) appendRow(chk *chunk.Chunk, row []interface{}) {
	for i, col := range row {
		if col == nil {
			chk.AppendNull(i)
			continue
		}
		switch x := col.(type) {
		case nil:
			chk.AppendNull(i)
		case int:
			chk.AppendInt64(i, int64(x))
		case int64:
			chk.AppendInt64(i, x)
		case uint64:
			chk.AppendUint64(i, x)
		case float64:
			chk.AppendFloat64(i, x)
		case float32:
			chk.AppendFloat32(i, x)
		case string:
			chk.AppendString(i, x)
		case []byte:
			chk.AppendBytes(i, x)
		case types.BinaryLiteral:
			chk.AppendBytes(i, x)
		case *types.MyDecimal:
			chk.AppendMyDecimal(i, x)
		case types.Time:
			chk.AppendTime(i, x)
		case json.BinaryJSON:
			chk.AppendJSON(i, x)
		case types.Duration:
			chk.AppendDuration(i, x)
		case types.Enum:
			chk.AppendEnum(i, x)
		case types.Set:
			chk.AppendSet(i, x)
		default:
			chk.AppendNull(i)
		}
	}
}
