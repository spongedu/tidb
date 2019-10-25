package log

import (
	"time"

	"github.com/pingcap/tidb/infoschema/inspection/log/item"
	"github.com/pingcap/tidb/infoschema/inspection/log/parser"
)

var (
	tidbLogPath = "/Users/duchuan.dc/pingcap/hackathon2019/tidb_log/"
	fw          []*parser.FileWrapper
)

type TiDBLogItem struct {
	Host      string         `json:"host"`
	Port      string         `json:"port"`
	Component string         `json:"component"`
	FileName  string         `json:"file"`
	Time      time.Time      `json:"time"`
	Level     item.LevelType `json:"level"`
	Content   string         `json:"content"`
}

type TiDBLogBatch struct {
	Logs []*TiDBLogItem `json:"logs"`
	Cnt  int            `json:"cnt"`
}

func ResetTiDBLogPath(p string) {
	tidbLogPath = p
}

func GetTiDBLogPath() string {
	return tidbLogPath
}
