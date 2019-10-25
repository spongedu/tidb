package main

import (
	"fmt"
	"time"

	"github.com/pingcap/tidb/infoschema/inspection/log"
	"github.com/pingcap/tidb/infoschema/inspection/log/search"
	"github.com/sirupsen/logrus"
)

func main() {
	//startTimeStr := "2019-10-24T11:35:29"
	  startTimeStr := "2019-10-24T11:35:29"
	//endTimeStr := "2019-10-24T11:35:47"
	  endTimeStr := "2019-10-24T11:35:47"


	TimeStampLayout := "2006-01-02T15:04:05"
	local, _ := time.LoadLocation("Asia/Chongqing")

	// startTime, _ := time.ParseInLocation(TimeStampLayout ,"1970-01-01T00:00:00", local)
	// endTime, _ := time.ParseInLocation(TimeStampLayout , "2030-01-01T00:00:00", local)
	startTime, err := time.ParseInLocation(TimeStampLayout ,startTimeStr, local)
	if err != nil {
		panic(err)
	}
	endTime, err := time.ParseInLocation(TimeStampLayout , endTimeStr, local)
	if err != nil {
		panic(err)
	}

	logrus.Infof("st=%s", startTime)
	logrus.Infof("et=%s", endTime)
	se, err := search.NewSequence(log.GetTiDBLogPath(), startTime, endTime)
	if err != nil {
		panic(err)
	}

	cnt := 0
	for {
		item, err := se.Next()
		if err != nil {
			fmt.Println(err)
			break
		}
		if item != nil {
			cnt = cnt + 1
			fmt.Printf("C=%s", string(item.GetContent()))
		} else {
			break
		}
	}
	fmt.Printf("cnt=%d",cnt)
}
