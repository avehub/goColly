package motor

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"spider/logic"
	"spider/models"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

func init() {
	rand.NewSource(time.Now().UnixNano())
}

type Coll struct {
	C *colly.Collector
}

const (
	// 摩托车网站地址
	DOMAIN = "https://motor.xcar.com.cn/select/"
	// 类型，摩托 = 2
	TYPE = 2
	// 设置最大重试次数
	MAX = 3
)

var wait sync.WaitGroup

func Start() {
	c := colly.NewCollector(
		// 设置浏览器UA头
		// colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
		// 设置网络异步请求
		// colly.Async(true),
		// 限制URL的递归深度
		colly.MaxDepth(2),
		// 忽略robots协议
		colly.IgnoreRobotsTxt(),
		// 设置GET请求本地缓存文件夹
		colly.CacheDir("cache-dir"),
	)
	Coll := Coll{C: c}
	Coll.Brand()
	// 请求前
	c.OnRequest(func(request *colly.Request) {
		// currentProxyIndex := rand.Intn(52)
		// request.ProxyURL = "http://" + logic.ProxyList[currentProxyIndex]
	})

	// 请求完成后回调
	c.OnResponse(func(response *colly.Response) {
		fmt.Println("完成回调前")

		// 网站反应渲染比较慢
		time.Sleep(time.Second * 3)
		fmt.Println("完成回调后")
	})

	//请求发生错误回调
	c.OnError(func(response *colly.Response, err error) {
		if response.StatusCode == http.StatusServiceUnavailable {
			for i := 0; i < MAX; i++ {
				time.Sleep(time.Second * 3)
				err := c.Visit(response.Request.URL.String())
				if err == nil {
					fmt.Println("Retry successfully")
					break
				} else {
					fmt.Println(err)
					log.Fatalln(err)
				}
			}
		}
		fmt.Printf("发生错误了:%v", err)
	})

	// 在所有OnHTML之后执行，可以用来做一些回收操作
	c.OnScraped(func(response *colly.Response) {
		fmt.Println("----> 所有匹配已完成")
	})

	// 开始访问
	err := c.Visit(DOMAIN)
	if err != nil {
		log.Fatalln(err)
		return
	}

}

// 品牌数据解析
func (Colly *Coll) Brand() {
	brand := models.Brand{}
	LetterMap := logic.LetterMap()
	// 因为DOM标签均在js内部，这块直接使用OnHTML会匹配不到，所以需要在OnScraped单独处理一下，或者直接改用访问已经加载好后的页面（本地源码）
	Colly.C.OnHTML(".list_cont", func(el *colly.HTMLElement) {
		el.ForEach("li", func(_ int, ei *colly.HTMLElement) {
			// 品牌首字母
			letter := ei.ChildText(".name")
			fmt.Printf("品牌首字母%v\n", letter)
			name := ei.ChildText("a")
			fmt.Printf("品牌名称：%v\n", name)
			logo := ei.ChildAttr("a > span > img", "src")
			if !strings.Contains(logo, "http") {
				logo = "https:" + logo
			}
			fmt.Printf("品牌Logo：%v\n", logo)
			// 品牌详情
			detail := ei.ChildAttr("a", "href")
			if !strings.Contains(detail, "http") {
				detail = "https:" + detail
			}
			brand.Letter = LetterMap[letter]
			brand.Logo = logo
			brand.Name = name
			brand.ID = 0
			brand.Type = TYPE
			fmt.Printf("%v\n", brand)
			// models.AddBrandOne(&brand)
			wait.Add(1)
			Coll := Coll{C: Colly.C}
			go Coll.Bank(&detail, &brand.ID)
			wait.Wait()

			fmt.Printf("\n")
		})

	})
}

// 系列数据解析
func (Colly *Coll) Bank(url *string, id *uint) {
	brand := models.Brand{}
	data := []models.Brand{}
	Colly.C.OnHTML(".BrandRight", func(ed *colly.HTMLElement) {
		childName := ed.ChildText(".MotoBrand > h3")
		fmt.Printf("子品牌名称：%v\n", childName)

		ed.ForEach(".MotoBrand > ul", func(_ int, ei *colly.HTMLElement) {
			bankName := ei.ChildText(".moto_column > .name > a")
			fmt.Printf("系列名称%v\n", bankName)
			bankImg := ei.ChildAttr(".MotoImg > img", "src")
			if !strings.Contains(bankImg, "http") {
				bankImg = "https:" + bankImg
			}
			fmt.Printf("系列照片：%v\n", bankImg)

			brand.Logo = bankImg
			brand.Name = bankName
			brand.Pid = *id
			brand.Type = TYPE
			data = append(data, brand)

			fmt.Printf("\n")
		})

		fmt.Printf("%v\n", data)
		// models.AddBrandBatch(&data)
	})

	wait.Done()
}
