package pholcus_lib

// 基础包
import (
	"github.com/molast/crawler-core/app/downloader/request" //必需
	. "github.com/molast/crawler-core/app/spider"           //必需
	. "github.com/molast/crawler-core/app/spider/common"    //选用
	"net/http"                                              //设置http.Header
)

func init() {
	Lewa.Register()
}

var Lewa = &Spider{
	Name:        "乐蛙登录测试",
	Description: "乐蛙登录测试 [Auto Page] [http://accounts.lewaos.com]",
	// Pausetime: 300,
	// Keyin:   KEYIN,
	// Limit:        LIMIT,
	EnableCookie: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.AddQueue(&request.Request{Url: "http://accounts.lewaos.com/", Rule: "登录页"})
		},

		Trunk: map[string]*Rule{

			"登录页": {
				ParseFunc: func(ctx *Context) {
					// ctx.AddQueue(&request.Request{
					// 	Url:    "http://accounts.lewaos.com",
					// 	Rule:   "登录后",
					// 	Method: "POST",
					// 	PostData: "username=123456@qq.com&password=123456&login_btn=login_btn&submit=login_btn",
					// })
					NewForm(
						ctx,
						"登录后",
						"http://accounts.lewaos.com",
						ctx.GetDom().Find(".userlogin.lw-pl40"),
					).Inputs(map[string]string{
						"username": "",
						"password": "",
					}).Submit()
				},
			},
			"登录后": {
				ParseFunc: func(ctx *Context) {
					// 结果存入Response中转
					ctx.Output(map[string]interface{}{
						"Body":   ctx.GetText(),
						"Cookie": ctx.GetCookie(),
					})
					ctx.AddQueue(&request.Request{
						Url:    "http://accounts.lewaos.com/member",
						Rule:   "个人中心",
						Header: http.Header{"Referer": []string{ctx.GetUrl()}},
					})
				},
			},
			"个人中心": {
				ParseFunc: func(ctx *Context) {
					// 结果存入Response中转
					ctx.Output(map[string]interface{}{
						"Body":   ctx.GetText(),
						"Cookie": ctx.GetCookie(),
					})
				},
			},
		},
	},
}
