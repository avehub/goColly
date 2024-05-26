package car

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func Common() {
	c := colly.NewCollector(
		// 设置浏览器UA头
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36 Edg/113.0.1774.42"),
		// 设置网络异步请求
		colly.Async(true),
		// 限制URL的递归深度
		colly.MaxDepth(2),
		// 忽略robots协议
		colly.IgnoreRobotsTxt(),
		// 设置GET请求本地缓存文件夹
		colly.CacheDir("cache-dir"),
	)

	// 请求前
	c.OnRequest(func(request *colly.Request) {
		// 设置浏览器UA头
		fmt.Println("User-Agent:", request.Headers.Get("User-Agent"))
	})

}
