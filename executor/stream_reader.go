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

	"github.com/pingcap/errors"

	//"github.com/pingcap/parser/model"
	//plannercore "github.com/pingcap/tidb/planner/core"
	//"github.com/pingcap/tidb/statistics"
	//"github.com/pingcap/tidb/table"
	"github.com/pingcap/tidb/util/chunk"
	//"github.com/pingcap/tidb/util/ranger"
	//tipb "github.com/pingcap/tipb/go-tipb"
	"golang.org/x/net/context"
)

// make sure `TableReaderExecutor` implements `Executor`.
var _ Executor = &StreamReaderExecutor{}

// StreamReaderExecutor sends DAG request and reads data from a stream.
type StreamReaderExecutor struct {
	baseExecutor

	/*
	table           table.Table
	physicalTableID int64
	keepOrder       bool
	desc            bool
	ranges          []*ranger.Range
	dagPB           *tipb.DAGRequest
	// columns are only required by union scan.
	columns []*model.ColumnInfo

	// resultHandler handles the order of the result. Since (MAXInt64, MAXUint64] stores before [0, MaxInt64] physically
	// for unsigned int.
	resultHandler *tableResultHandler
	streaming     bool
	feedback      *statistics.QueryFeedback

	// corColInFilter tells whether there's correlated column in filter.
	corColInFilter bool
	// corColInAccess tells whether there's correlated column in access conditions.
	corColInAccess bool
	plans          []plannercore.PhysicalPlan
	*/
	rowCnt		  int
}

// Open initialzes necessary variables for using this executor.
func (e *StreamReaderExecutor) Open(ctx context.Context) error {
	return nil
}

// Next fills data into the chunk passed by its caller.
func (e *StreamReaderExecutor) Next(ctx context.Context, chk *chunk.Chunk) error {
	chk.Reset()
	for {
		if e.rowCnt < 20 && chk.NumRows() < e.ctx.GetSessionVars().MaxChunkSize {
			chk.AppendInt64(0, int64(e.rowCnt))
			e.rowCnt += 1
		} else {
			break
		}
	}
	return errors.Trace(nil)
}

// Close implements the Executor Close interface.
func (e *StreamReaderExecutor) Close() error {
	return nil
}

