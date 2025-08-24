package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/molast/crawler-core/app"
	"github.com/molast/crawler-core/runtime/cache"
	"github.com/molast/crawler-core/runtime/status"

	_ "github.com/molast/crawler-core/pholcus_lib"
)

type CrawlerRequest struct {
	Keyins         string `json:"keyins"`          // 自定义输入，后期切分为多个任务的Keyin自定义配置
	Limit          int64  `json:"limit"`           // 采集上限，0为不限，若在规则中设置初始值为LIMIT则为自定义限制，否则默认限制请求数
	OutType        string `json:"out_type"`        // 输出方式 (如 json/csv/db)
	ThreadNum      int    `json:"thread_num"`      // 全局最大并发量
	Pausetime      int64  `json:"pause_time"`      // 暂停时长参考/ms(随机: Pausetime/2 ~ Pausetime*2)
	ProxySecond    int64  `json:"proxy_second"`    // 代理IP更换的间隔秒数
	DockerCap      int    `json:"docker_cap"`      // 分段转储容器容量
	SuccessInherit bool   `json:"success_inherit"` // 继承历史成功记录
	FailureInherit bool   `json:"failure_inherit"` // 继承历史失败记录
}

func main() {
	app.LogicApp.Init(cache.Task.Mode, cache.Task.Port, cache.Task.Master)
	if cache.Task.Mode == status.UNSET {
		return
	}

	go func() {
		run()
	}()

	http.HandleFunc("/add_task", addTaskHandler)
	log.Println("HTTP 服务已启动: http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}
	var req CrawlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// 这里可以保存到全局 cache 或者任务队列
	// 比如：cache.Task = req

	log.Printf("收到新任务: %+v\n", req)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"task":   req,
	})
}

func run() {
	sps := app.LogicApp.GetSpiderLib()
	app.LogicApp.SpiderPrepare(sps).Run()
}
