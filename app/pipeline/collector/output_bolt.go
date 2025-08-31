package collector

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/molast/crawler-core/app/pipeline/collector/data"
	"github.com/molast/crawler-core/common/buckets"
	"github.com/molast/crawler-core/common/util"
	"github.com/molast/crawler-core/config"
	"github.com/molast/crawler-core/logs"
)

/************************ boltDB 输出 ***************************/
func init() {
	DataOutput["bolt"] = func(self *Collector) (err error) {
		var (
			namespace = util.FileNameReplace(self.namespace())
		)
		defer func() {
			if p := recover(); p != nil {
				err = fmt.Errorf("%v", p)
			}
		}()
		folder := config.FILE_DIR + "/" + "crawler.db"
		bx, err := buckets.Open(folder)
		if err != nil {
			logs.Log.Error(fmt.Sprintf("open bucket failed: %v", err))
			return err
		}
		defer bx.Close()

		bucket, err := bx.New([]byte(namespace))
		if err != nil {
			logs.Log.Error(fmt.Sprintf("couldn't create todos bucket: %v", err))
			return err
		}
		dataKey := self.GetKeyin()
		if dataKey == "" {
			return fmt.Errorf("dataKey is empty")
		}

		list := make([]data.DataCell, len(self.dataDocker))
		for i, cell := range self.dataDocker {
			for k, v := range cell["Data"].(map[string]interface{}) {
				cell[k] = v
			}
			delete(cell, "Data")
			delete(cell, "RuleName")
			if !self.Spider.OutDefaultField() {
				delete(cell, "Url")
				delete(cell, "ParentUrl")
				delete(cell, "DownloadTime")
			}
			list[i] = cell
		}

		bs, err := jsoniter.Marshal(list)
		if err != nil {
			return err
		}

		items := []struct {
			Key   []byte
			Value []byte
		}{
			{[]byte(dataKey), bs},
		}

		err = bucket.Insert(items)
		if err != nil {
			return err
		}
		return
	}
}
