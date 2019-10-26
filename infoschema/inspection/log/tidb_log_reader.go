package log

import (
	"path"
	"time"

	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/infoschema/inspection/log/item"
	"github.com/pingcap/tidb/infoschema/inspection/log/parser"
)

var (
	tidbLogPath = "/Users/duchuan.dc/pingcap/hackathon2019/tidb_log/"
	fw          []*parser.FileWrapper
)

type TiDBLogItem struct {
	Address      string      `json:"address"`
	FileName  string         `json:"file"`
	Time      time.Time      `json:"time"`
	Level     string 		 `json:"level"`
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
	dir, _ := path.Split(config.GetGlobalConfig().Log.File.Filename)
	return dir
}

func ParseLevelToStr(level item.LevelType) string {
	if level == item.LevelInvalid {
		return "invalid"
	}
	if level == item.LevelFATAL {
		return "fatal"
	}
	if level == item.LevelERROR {
		return "error"
	}
	if level == item.LevelWARN {
		return "warn"
	}
	if level == item.LevelINFO {
		return "info"
	}
	if level == item.LevelDEBUG {
		return "debug"
	}
	return "unknown"
}
