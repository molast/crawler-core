package model

type CrawlerRequest struct {
	Spiders        []string `json:"spiders"`         // 需要执行任务的名称
	Keyins         string   `json:"keyins"`          // 自定义配置，后期切分为多个任务的Keyin自定义配置
	ThreadNum      int      `json:"thread_num"`      // 并发协程
	Limit          int64    `json:"limit"`           // 采集上限，0为不限，若在规则中设置初始值为LIMIT则为自定义限制，否则默认限制请求数
	DockerCap      int      `json:"docker_cap"`      // 分批输出限制
	Pausetime      int64    `json:"pause_time"`      // 暂停时长参考 /ms(随机: Pausetime/2 ~ Pausetime*2)
	ProxySecond    int64    `json:"proxy_second"`    // 代理IP更换频率，间隔秒数
	OutType        string   `json:"out_type"`        // 输出方式 (csv/excel/kafka/mgo/mysql)
	SuccessInherit bool     `json:"success_inherit"` // 继承并保存成功记录
	FailureInherit bool     `json:"failure_inherit"` // 继承并保存失败记录
}
