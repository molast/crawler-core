package pipeline

import (
	"sort"

	"github.com/molast/crawler-core/app/pipeline/collector"
	"github.com/molast/crawler-core/common/kafka"
	"github.com/molast/crawler-core/common/mgo"
	"github.com/molast/crawler-core/common/mysql"
	"github.com/molast/crawler-core/runtime/cache"
)

// 初始化输出方式列表collector.DataOutputLib
func init() {
	for out, _ := range collector.DataOutput {
		collector.DataOutputLib = append(collector.DataOutputLib, out)
	}
	sort.Strings(collector.DataOutputLib)
}

// RefreshOutput 刷新输出方式的状态
func RefreshOutput() {
	switch cache.Task.OutType {
	case "mgo":
		mgo.Refresh()
	case "mysql":
		mysql.Refresh()
	case "kafka":
		kafka.Refresh()
	}
}
