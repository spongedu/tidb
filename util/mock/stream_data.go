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

package mock

import (
	"encoding/json"
	"fmt"

	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/tidb/types"
	log "github.com/sirupsen/logrus"
)

// Mock data for stream reader.
var MockStreamData []Event

type Event struct {
	ID         int64      `json:"id"`
	Content    string     `json:"content"`
	CreateTime types.Time `json:"create_time"`
}

func init() {
	t := types.CurrentTime(mysql.TypeTimestamp)

	for i := 0; i < 10000; i++ {
		tt := types.Time{
			Time: types.FromDate(t.Time.Year(), t.Time.Month(), t.Time.Day(), t.Time.Hour(), t.Time.Minute(), t.Time.Second()+i, t.Time.Microsecond()),
			Type: mysql.TypeTimestamp,
			Fsp:  types.DefaultFsp,
		}
		content := fmt.Sprintf("TiDB Stream Test Data %d", i)
		evt := Event{int64(i), content, tt}
		MockStreamData = append(MockStreamData, evt)
	}

	for i := 0; i < 10; i++ {
		data, err := json.Marshal(MockStreamData[i])
		if err != nil {
			log.Fatalf("[mock stream data marshal failed]%v", err)
		}

		evt := Event{}
		err = json.Unmarshal([]byte(data), &evt)
		if err != nil {
			log.Fatalf("[mock stream data unmarshal failed]%v", err)
		}

		MockStreamData[i] = evt
		log.Errorf("[mock stream data]%v-%s", MockStreamData[i], string(data))
	}
}
