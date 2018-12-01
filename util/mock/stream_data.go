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

// Mock data for stream reader demo.
var (
	MockStreamData       []Event
	MockKafkaStreamData  []KafkaEvent
	MockPulsarStreamData []PulsarEvent
	MockLogStreamData    []LogEvent
)

type Event struct {
	ID         int64      `json:"id"`
	Content    string     `json:"content"`
	CreateTime types.Time `json:"create_time"`
}

type KafkaEvent struct {
	ID         int64      `json:"id"`
	Content    string     `json:"content"`
	CreateTime types.Time `json:"create_time"`
}

type PulsarEvent struct {
	ID         int64      `json:"id"`
	Content    string     `json:"content"`
	CreateTime types.Time `json:"create_time"`
}

type LogEvent struct {
	ID         int64      `json:"id"`
	Content    string     `json:"content"`
	CreateTime types.Time `json:"create_time"`
}

func genData(i int, t types.Time) Event {
	tt := types.Time{
		Time: types.FromDate(t.Time.Year(), t.Time.Month(), t.Time.Day(), t.Time.Hour(), t.Time.Minute(), t.Time.Second()+i, t.Time.Microsecond()),
		Type: mysql.TypeTimestamp,
		Fsp:  types.DefaultFsp,
	}
	content := fmt.Sprintf("TiDB Stream Demo Data %d", i)
	evt := Event{int64(i), content, tt}

	return evt
}

func genKafkaData(i int, t types.Time) KafkaEvent {
	tt := types.Time{
		Time: types.FromDate(t.Time.Year(), t.Time.Month(), t.Time.Day(), t.Time.Hour(), t.Time.Minute(), t.Time.Second()+i, t.Time.Microsecond()),
		Type: mysql.TypeTimestamp,
		Fsp:  types.DefaultFsp,
	}
	content := fmt.Sprintf("TiDB Stream Demo Data %d", i)
	evt := KafkaEvent{int64(i), content, tt}

	return evt
}

func genPulsarData(i int, t types.Time) PulsarEvent {
	evt := PulsarEvent{}
	return evt
}

func genLogData(i int, t types.Time) LogEvent {
	evt := LogEvent{}
	return evt
}

func init() {
	t := types.CurrentTime(mysql.TypeTimestamp)

	for i := 0; i < 100; i++ {
		MockStreamData = append(MockStreamData, genData(i, t))
		MockKafkaStreamData = append(MockKafkaStreamData, genKafkaData(i, t))
		MockPulsarStreamData = append(MockPulsarStreamData, genPulsarData(i, t))
		MockLogStreamData = append(MockLogStreamData, genLogData(i, t))
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

		log.Errorf("[mock stream kafka data]%v", MockKafkaStreamData[i])
		log.Errorf("[mock stream pulsar data]%v", MockPulsarStreamData[i])
		log.Errorf("[mock stream log data]%v", MockLogStreamData[i])
	}
}
