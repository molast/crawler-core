package pholcus_lib

// 基础包
import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/molast/crawler-core/app/downloader/request"
	. "github.com/molast/crawler-core/app/spider"
)

func init() {
	Wangyi.Register()
}

var Wangyi = &Spider{
	Name:        "网易国内新闻",
	Description: "网易国内新闻",
	// Pausetime:    300,
	// Keyin:        KEYIN,
	Limit:        50,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.AddQueue(&request.Request{Url: "http://news.163.com/domestic",
				Rule:   "新闻列表",
				Header: http.Header{"Content-Type": []string{"application/xml"}},
			})
		},

		Trunk: map[string]*Rule{
			"新闻列表": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					//获取新闻列表
					newsList := query.Find(".ndi_main")
					fmt.Println(newsList.Html())
					newsList.Each(func(i int, s *goquery.Selection) {
						imgUrl := ""
						// 预览图
						if url, ok := s.Find("a img").Attr("src"); ok {
							imgUrl = url
						}
						titleA := s.Find(".news_title a")
						// 标题
						newsTitle := titleA.Text()
						// 时间
						newsTime := s.Find(".news_tag").Text()

						if url, ok := titleA.Attr("href"); ok {
							ctx.AddQueue(&request.Request{
								Url:  url,
								Rule: "新闻内容",
								Temp: map[string]interface{}{
									"imgUrl":    imgUrl,
									"newsTitle": newsTitle,
									"newsTime":  newsTime,
								},
							})
						}
					})
				},
			},

			"新闻内容": {
				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"标题",
					"内容",
					"预览图",
					"发布时间",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					//// 若有多页内容，则获取阅读全文的链接并获取内容
					//if pageAll := query.Find(".post_body"); len(pageAll.Nodes) != 0 {
					//	if pageAllUrl, ok := pageAll.Attr("href"); ok {
					//		ctx.AddQueue(&request.Request{
					//			Url:  pageAllUrl,
					//			Rule: "热点新闻",
					//			Temp: ctx.CopyTemps(),
					//		})
					//	}
					//	return
					//}

					// 获取内容
					content, _ := query.Find(".post_body").Html()

					re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
					// content = re.ReplaceAllStringFunc(content, strings.ToLower)
					content = re.ReplaceAllString(content, "")

					// 获取发布日期
					release := query.Find(".post_info").Text()
					release = strings.Split(release, "来源:")[0]
					release = strings.Trim(release, " \t\n")

					// 结果存入Response中转
					ctx.Output(map[int]interface{}{
						0: ctx.GetTemp("newsTitle", ""),
						1: content,
						2: ctx.GetTemp("imgUrl", ""),
						4: release,
					})
				},
			},
		},
	},
}
