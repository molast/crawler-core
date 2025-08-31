package config

import (
	"strings"

	"github.com/molast/crawler-core/logs/logs"
	"github.com/molast/crawler-core/runtime/cache"
)

const (
	TAG string = "crawler" // 软件标识符
)

// 默认配置。
const (
	WORK_ROOT             = TAG + "_temp"                   // 运行时的目录名称
	CONFIG                = WORK_ROOT + "/config.yaml"      // 配置文件路径
	CACHE_DIR             = WORK_ROOT + "/cache"            // 缓存文件目录
	LOG                   = WORK_ROOT + "/logs/crawler.log" // 日志文件路径
	LOG_ASYNC      bool   = true                            // 是否异步输出日志
	PHANTOMJS_TEMP        = CACHE_DIR                       // Surfer-Phantom下载器：js文件临时目录
	HISTORY_TAG    string = "history"                       // 历史记录的标识符
	HISTORY_DIR           = WORK_ROOT + "/" + HISTORY_TAG   // excel或csv输出方式下，历史记录目录
	SPIDER_EXT     string = ".crawler.html"                 // 动态规则扩展名
)

// 来自配置文件的配置项。
var (
	CRAWLS_CAP               = setting.GetInt("crawlcap")      // 蜘蛛池最大容量
	PHANTOMJS                = setting.GetString("phantomjs")  // Surfer-Phantom下载器：phantomjs程序路径
	PROXY                    = setting.GetString("proxylib")   // 代理IP文件路径
	SPIDER_DIR               = setting.GetString("spiderdir")  // 动态规则目录
	FILE_DIR                 = setting.GetString("fileoutdir") // 文件（图片、HTML等）结果的输出目录
	TEXT_DIR                 = setting.GetString("textoutdir") // excel或csv输出方式下，文本结果的输出目录
	DB_NAME                  = setting.GetString("dbname")     // 数据库名称
	MGO_ADMIN_USERNAME       = setting.GetString("mgo.username")
	MGO_ADMIN_PASSWORD       = setting.GetString("mgo.password")
	MGO_CONN_STR             = setting.GetString("mgo.connstring")              // mongodb连接字符串
	MGO_CONN_CAP             = setting.GetInt("mgo.conncap")                    // mongodb连接池容量
	MGO_CONN_GC_SECOND       = setting.GetInt64("mgo.conngcsecond")             // mongodb连接池GC时间，单位秒
	MYSQL_CONN_STR           = setting.GetString("mysql.connstring")            // mysql连接字符串
	MYSQL_CONN_CAP           = setting.GetInt("mysql.conncap")                  // mysql连接池容量
	MYSQL_MAX_ALLOWED_PACKET = setting.GetInt("mysql.maxallowedpacket")         // mysql通信缓冲区的最大长度
	KAFKA_BORKERS            = setting.GetString("kafka.brokers")               // kafka brokers
	LOG_CAP                  = setting.GetInt64("log.cap")                      // 日志缓存的容量
	LOG_LEVEL                = logLevel(setting.GetString("log.level"))         // 全局日志打印级别（亦是日志文件输出级别）
	LOG_CONSOLE_LEVEL        = logLevel(setting.GetString("log.consolelevel"))  // 日志在控制台的显示级别
	LOG_FEEDBACK_LEVEL       = logLevel(setting.GetString("log.feedbacklevel")) // 客户端反馈至服务端的日志级别
	LOG_LINEINFO             = setting.GetBool("log.lineinfo")                  // 日志是否打印行信息                                  // 客户端反馈至服务端的日志级别
	LOG_SAVE                 = setting.GetBool("log.save")                      // 是否保存所有日志到本地文件
)

func init() {
	// 主要运行时参数的初始化
	cache.Task = &cache.AppConf{
		Mode:            setting.GetInt("run.mode"),             // 节点角色
		Port:            setting.GetInt("run.port"),             // 主节点端口
		Master:          setting.GetString("run.master"),        // 服务器(主节点)地址，不含端口
		ThreadNum:       setting.GetInt("run.thread"),           // 全局最大并发量
		Pausetime:       setting.GetInt64("run.pause"),          // 暂停时长参考/ms(随机: Pausetime/2 ~ Pausetime*2)
		OutType:         setting.GetString("run.outtype"),       // 输出方式
		DockerCap:       setting.GetInt("run.dockercap"),        // 分段转储容器容量
		Limit:           setting.GetInt64("run.limit"),          // 采集上限，0为不限，若在规则中设置初始值为LIMIT则为自定义限制，否则默认限制请求数
		ProxySecond:     setting.GetInt64("run.proxysecond"),    // 代理IP更换的间隔秒钟数
		SuccessInherit:  setting.GetBool("run.success"),         // 继承历史成功记录
		FailureInherit:  setting.GetBool("run.failure"),         // 继承历史失败记录
		AutoOpenBrowser: setting.GetBool("run.autoopenbrowser"), // 自动打开浏览器
	}
}

func logLevel(l string) int {
	switch strings.ToLower(l) {
	case "app":
		return logs.LevelApp
	case "emergency":
		return logs.LevelEmergency
	case "alert":
		return logs.LevelAlert
	case "critical":
		return logs.LevelCritical
	case "error":
		return logs.LevelError
	case "warning":
		return logs.LevelWarning
	case "notice":
		return logs.LevelNotice
	case "informational":
		return logs.LevelInformational
	case "info":
		return logs.LevelInformational
	case "debug":
		return logs.LevelDebug
	}
	return -10
}
