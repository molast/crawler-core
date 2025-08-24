package config

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/molast/crawler-core/runtime/status"
	"github.com/spf13/viper"
)

// 配置文件涉及的默认配置。
const (
	crawlcap              int    = 50                                          // 蜘蛛池最大容量
	datachancap           int    = 2 << 14                                     // 收集器容量(默认65536)
	logcap                int64  = 10000                                       // 日志缓存的容量
	loglevel              string = "debug"                                     // 全局日志打印级别（亦是日志文件输出级别）
	logconsolelevel       string = "info"                                      // 日志在控制台的显示级别
	logfeedbacklevel      string = "error"                                     // 客户端反馈至服务端的日志级别
	loglineinfo           bool   = false                                       // 日志是否打印行信息
	logsave               bool   = true                                        // 是否保存所有日志到本地文件
	phantomjs                    = WORK_ROOT + "/phantomjs"                    // phantomjs文件路径
	proxylib              string = "需手动输入，当前仅支持http://www.goubanjia.com/提供的链接" // 代理ip商提供的地址
	spiderdir                    = WORK_ROOT + "/spiders"                      // 动态规则目录
	fileoutdir                   = WORK_ROOT + "/file_out"                     // 文件（图片、HTML等）结果的输出目录
	textoutdir                   = WORK_ROOT + "/text_out"                     // excel或csv输出方式下，文本结果的输出目录
	dbname                       = TAG                                         // 数据库名称
	mgoconnstring         string = "127.0.0.1:27017"                           // mongodb连接字符串
	mgoconncap            int    = 1024                                        // mongodb连接池容量
	mgoconngcsecond       int64  = 600                                         // mongodb连接池GC时间，单位秒
	mysqlconnstring       string = "root:@tcp(127.0.0.1:3306)"                 // mysql连接字符串
	mysqlconncap          int    = 2048                                        // mysql连接池容量
	mysqlmaxallowedpacket int    = 1048576                                     //mysql通信缓冲区的最大长度，单位B，默认1MB
	kafkabrokers          string = "127.0.0.1:9092"                            //kafka broker字符串,逗号分割

	mode                   = status.UNSET // 节点角色
	autoOpenBrowser bool   = false        // 是否自动打开浏览器
	port            int    = 2015         // 主节点端口
	master          string = "127.0.0.1"  // 服务器(主节点)地址，不含端口
	thread          int    = 20           // 全局最大并发量
	pause           int64  = 300          // 暂停时长参考/ms(随机: Pausetime/2 ~ Pausetime*2)
	outtype         string = "csv"        // 输出方式
	dockercap       int    = 10000        // 分段转储容器容量
	limit           int64  = 0            // 采集上限，0为不限，若在规则中设置初始值为LIMIT则为自定义限制，否则默认限制请求数
	proxysecond     int64  = 0            // 代理IP更换的间隔秒钟数
	success         bool   = true         // 继承历史成功记录
	failure         bool   = true         // 继承历史失败记录
)

var (
	_viper *viper.Viper
	_once  sync.Once
)

var setting = initViper()

func initViper() *viper.Viper {
	_once.Do(func() {
		v := viper.New()
		v.SetConfigFile(CONFIG)
		v.SetConfigType("yaml")

		// 先确保目录存在
		_ = os.MkdirAll(filepath.Clean(HISTORY_DIR), 0777)
		_ = os.MkdirAll(filepath.Clean(CACHE_DIR), 0777)
		_ = os.MkdirAll(filepath.Clean(PHANTOMJS_TEMP), 0777)

		// 尝试读取配置文件
		if err := v.ReadInConfig(); err != nil {
			// 文件不存在，写入默认配置
			defaultConfig(v)
			if err := v.WriteConfigAs(CONFIG); err != nil {
				log.Fatalf("写入默认配置失败: %v", err)
			}
		} else {
			// 文件存在，检查并填充缺省值
			trySet(v)
			if err := v.WriteConfig(); err != nil {
				log.Fatalf("更新配置失败: %v", err)
			}
		}

		// 创建配置中的目录
		_ = os.MkdirAll(filepath.Clean(v.GetString("spiderdir")), 0777)
		_ = os.MkdirAll(filepath.Clean(v.GetString("fileoutdir")), 0777)
		_ = os.MkdirAll(filepath.Clean(v.GetString("textoutdir")), 0777)
		_viper = v
	})
	return _viper
}

func defaultConfig(v *viper.Viper) {
	v.SetDefault("crawlcap", crawlcap)
	v.SetDefault("log.cap", logcap)
	v.SetDefault("log.level", loglevel)
	v.SetDefault("log.consolelevel", logconsolelevel)
	v.SetDefault("log.feedbacklevel", logfeedbacklevel)
	v.SetDefault("log.lineinfo", loglineinfo)
	v.SetDefault("log.save", logsave)
	v.SetDefault("phantomjs", phantomjs)
	v.SetDefault("proxylib", proxylib)
	v.SetDefault("spiderdir", spiderdir)
	v.SetDefault("fileoutdir", fileoutdir)
	v.SetDefault("textoutdir", textoutdir)
	v.SetDefault("dbname", dbname)
	v.SetDefault("mgo.username", "")
	v.SetDefault("mgo.password", "")
	v.SetDefault("mgo.connstring", mgoconnstring)
	v.SetDefault("mgo.conngcsecond", mgoconngcsecond)
	v.SetDefault("mgo.conncap", mgoconncap)
	v.SetDefault("mysql.connstring", mysqlconnstring)
	v.SetDefault("mysql.conncap", mysqlconncap)
	v.SetDefault("mysql.maxallowedpacket", mysqlmaxallowedpacket)
	v.SetDefault("kafka.brokers", kafkabrokers)
	v.SetDefault("run.mode", mode)
	v.SetDefault("run.port", port)
	v.SetDefault("run.master", master)
	v.SetDefault("run.thread", thread)
	v.SetDefault("run.pause", pause)
	v.SetDefault("run.outtype", outtype)
	v.SetDefault("run.dockercap", dockercap)
	v.SetDefault("run.limit", limit)
	v.SetDefault("run.proxysecond", proxysecond)
	v.SetDefault("run.success", success)
	v.SetDefault("run.failure", failure)
	v.SetDefault("run.autoopenbrowser", autoOpenBrowser)
}

func trySet(v *viper.Viper) {
	// crawlcap
	if v.GetInt("crawlcap") <= 0 {
		v.Set("crawlcap", crawlcap)
	}

	if v.GetInt("datachancap") <= 0 {
		v.Set("datachancap", datachancap)
	}

	// log 部分
	if v.GetInt64("log.cap") <= 0 {
		v.Set("log.cap", logcap)
	}
	level := v.GetString("log.level")
	if logLevel(level) == -10 {
		level = loglevel
	}
	v.Set("log.level", level)

	consoleLevel := v.GetString("log.consolelevel")
	if logLevel(consoleLevel) == -10 {
		consoleLevel = logconsolelevel
	}
	v.Set("log.consolelevel", logLevel2(consoleLevel, level))

	feedbackLevel := v.GetString("log.feedbacklevel")
	if logLevel(feedbackLevel) == -10 {
		feedbackLevel = logfeedbacklevel
	}
	v.Set("log.feedbacklevel", logLevel2(feedbackLevel, level))

	if !v.IsSet("log.lineinfo") {
		v.Set("log.lineinfo", loglineinfo)
	}
	if !v.IsSet("log.save") {
		v.Set("log.save", logsave)
	}

	// 路径/文件配置
	if v.GetString("phantomjs") == "" {
		v.Set("phantomjs", phantomjs)
	}
	if v.GetString("proxylib") == "" {
		v.Set("proxylib", proxylib)
	}
	if v.GetString("spiderdir") == "" {
		v.Set("spiderdir", spiderdir)
	}
	if v.GetString("fileoutdir") == "" {
		v.Set("fileoutdir", fileoutdir)
	}
	if v.GetString("textoutdir") == "" {
		v.Set("textoutdir", textoutdir)
	}
	if v.GetString("dbname") == "" {
		v.Set("dbname", dbname)
	}

	// mgo
	if v.GetString("mgo.connstring") == "" {
		v.Set("mgo.connstring", mgoconnstring)
	}
	if v.GetInt("mgo.conncap") <= 0 {
		v.Set("mgo.conncap", mgoconncap)
	}
	if v.GetInt64("mgo.conngcsecond") <= 0 {
		v.Set("mgo.conngcsecond", mgoconngcsecond)
	}

	// mysql
	if v.GetString("mysql.connstring") == "" {
		v.Set("mysql.connstring", mysqlconnstring)
	}
	if v.GetInt("mysql.conncap") <= 0 {
		v.Set("mysql.conncap", mysqlconncap)
	}
	if v.GetInt("mysql.maxallowedpacket") <= 0 {
		v.Set("mysql.maxallowedpacket", mysqlmaxallowedpacket)
	}

	// kafka
	if v.GetString("kafka.brokers") == "" {
		v.Set("kafka.brokers", kafkabrokers)
	}

	// run
	if v.GetInt("run.mode") < status.UNSET || v.GetInt("run.mode") > status.CLIENT {
		v.Set("run.mode", mode)
	}
	if v.GetInt("run.port") <= 0 {
		v.Set("run.port", port)
	}
	if v.GetString("run.master") == "" {
		v.Set("run.master", master)
	}
	if v.GetInt("run.thread") <= 0 {
		v.Set("run.thread", thread)
	}
	if v.GetInt64("run.pause") < 0 {
		v.Set("run.pause", pause)
	}
	if v.GetString("run.outtype") == "" {
		v.Set("run.outtype", outtype)
	}
	if v.GetInt("run.dockercap") <= 0 {
		v.Set("run.dockercap", dockercap)
	}
	if v.GetInt64("run.limit") < 0 {
		v.Set("run.limit", limit)
	}
	if v.GetInt64("run.proxysecond") <= 0 {
		v.Set("run.proxysecond", proxysecond)
	}
	if !v.IsSet("run.success") {
		v.Set("run.success", success)
	}
	if !v.IsSet("run.failure") {
		v.Set("run.failure", failure)
	}
	if !v.IsSet("run.autoopenbrowser") {
		v.Set("run.autoopenbrowser", autoOpenBrowser)
	}

	// 最后保存配置
	if err := v.WriteConfig(); err != nil {
		// 如果没有文件，用 WriteConfigAs 创建
		if err := v.WriteConfigAs(CONFIG); err != nil {
			log.Fatal("写配置失败:", err)
		}
	}
}

func logLevel2(l string, g string) string {
	a, b := logLevel(l), logLevel(g)
	if a < b {
		return l
	}
	return g
}
