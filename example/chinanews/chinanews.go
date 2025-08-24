package pholcus_lib

// 基础包
import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/molast/crawler-core/app/downloader/request" //必需
	. "github.com/molast/crawler-core/app/spider"           //必需
)

func init() {
	FileTest.Register()
}

var FileTest = &Spider{
	NotDefaultField: true, // 不添加默认字段
	Name:            "中国新闻网",
	Description:     "测试 [http://www.chinanews.com/scroll-news/news1.html]",
	// Pausetime: 300,
	// Keyin:   KEYIN,
	// Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.AddQueue(&request.Request{
				Url:  "http://www.chinanews.com/scroll-news/news1.html",
				Rule: "滚动新闻",
			})
		},

		Trunk: map[string]*Rule{
			"滚动新闻": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					//获取分页导航
					navBox := query.Find(".pagebox a")
					navBox.Each(func(i int, s *goquery.Selection) {
						if url, ok := s.Attr("href"); ok {
							ctx.AddQueue(&request.Request{
								Url:  "http://www.chinanews.com" + url,
								Rule: "新闻列表",
							})
						}
					})
				},
			},

			"新闻列表": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					//获取新闻列表
					newList := query.Find(".content_list li")
					newList.Each(func(i int, s *goquery.Selection) {
						// 新闻类型
						newsType := s.Find(".dd_lm a").Text()
						if newsType == "图片" {
							return
						}
						// 标题
						newsTitle := s.Find(".dd_bt a").Text()
						// 时间
						newsTime := s.Find(".dd_time").Text()
						if url, ok := s.Find(".dd_bt a").Attr("href"); ok {
							ctx.AddQueue(&request.Request{
								Url:  "http://www.chinanews.com" + url,
								Rule: "新闻内容",
								Temp: map[string]interface{}{
									"newsType":  newsType,
									"newsTitle": newsTitle,
									"newsTime":  newsTime,
								},
							})
						}
					})
				},
			},

			"新闻内容": {
				ItemFields: []string{
					"type",
					"origin",
					"title",
					"content",
					"time",
					"url",
					"parentUrl",
				},

				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					newsType := ctx.GetTemp("newsType", "")
					// 正文
					var content = ""
					if newsType == "视频" {
						query.Find(".content_desc").Children().Not(".content_editor").EachWithBreak(func(_ int, s *goquery.Selection) bool {
							htmlContent, _ := s.Html()
							content += fmt.Sprintf("%s\n", htmlContent) // 添加换行符
							return true
						})
					} else {
						content, _ = query.Find(".left_zw").Html()
					}

					// 来源
					from := strings.TrimSpace(query.Find(".content_left_time a").Text())
					if from == "" {
						s := ""
						if newsType == "视频" {
							s = strings.TrimSpace(query.Find(".content_title .left p").Text())
						} else {
							s = strings.TrimSpace(query.Find(".content_left_time").Contents().First().Text())
						}
						sep := "来源："
						i := strings.LastIndex(s, sep)
						//来源字符串特殊处理
						if i == -1 {
							from = "未知"
						} else {
							from = s[i+len(sep):]
						}
					}
					//输出格式
					ctx.Output(map[int]interface{}{
						0: newsType,
						1: from,
						2: ctx.GetTemp("newsTitle", ""),
						3: content,
						4: ctx.GetTemp("newsTime", ""),
						5: ctx.GetUrl(),
						6: ctx.GetReferer(),
					})
				},
			},
		},
	},
}
