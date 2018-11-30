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
	"strconv"
	"time"
)

// Mock data for stream reader.
var MockStreamData []Event

type Event struct {
	ID         int64  `json:"id"`
	Content    string `json:"name"`
	CreateTime int64  `json:"create_time"`
}

func init() {
	for i := 0; i < 10000; i++ {
		evt := Event{int64(i), strconv.Itoa(i), time.Now().Unix()}
		MockStreamData = append(MockStreamData, evt)
	}
}
