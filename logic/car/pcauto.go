package car

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

const (
	// 类型，汽车 = 1
	TYPE = 1
	// 设置最大重试次数
	MAX     = 3
	GODEBUG = 0
)

var wait sync.WaitGroup

func Start() {
	// 创建收集器
	c := colly.NewCollector()
	// 初始化
	// data := []models.Brand{}
	brand := models.Brand{}

	LetterMap := logic.LetterMap()
	// 品牌列表
	c.OnHTML(".brand", func(e *colly.HTMLElement) {
		// var i uint = 0
		e.ForEach(".js-letter", func(_ int, ed *colly.HTMLElement) {
			// 获取字母
			abc := ed.ChildText(".brand-index")
			fmt.Printf("字母序号：%v\n", abc)

			ed.ForEachWithBreak(".brand-item", func(j int, et *colly.HTMLElement) bool {
				// 获取品牌名称
				name := et.ChildText(".brand-name")
				fmt.Println(name)
				// 获取品牌Logo
				logo := et.ChildAttr(".brand-item-flex > .brand-logo > img", "src")
				fmt.Println(logo)
				// 品牌车型详情
				href := et.ChildAttr("a", "href")
				if !strings.Contains(href, "http") {
					href = "https:" + href
				}
				fmt.Println(href)

				brand.Letter = LetterMap[abc]
				brand.Logo = logo
				brand.Name = name
				brand.ID = 0
				brand.Type = TYPE
				models.AddBrandOne(&brand)
				wait.Add(1)
				go Brand(&href, &brand.ID)
				wait.Wait()
				// data = append(data, childData)

				return true
			})
			fmt.Printf("\n\n\n")
		})

		// 选择当前元素层级所有品牌
		// iocn := e.ChildText("")
		// models.AddBrandBatch(&data)
	})

	// 请求发起时回调,一般用来设置请求头等
	c.OnRequest(func(request *colly.Request) {
		fmt.Println("----> 开始请求了")
	})

	// 请求完成后回调
	c.OnResponse(func(response *colly.Response) {
		fmt.Println("----> 开始返回了")
	})

	//请求发生错误回调
	c.OnError(func(response *colly.Response, err error) {
		fmt.Printf("发生错误了:%v", err)
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

	})

	// 在所有OnHTML之后执行，可以用来做一些回收操作
	c.OnScraped(func(response *colly.Response) {
		fmt.Println("----> 所有匹配已完成")
	})

	// 开始访问
	err := c.Visit("https://price.pcauto.com.cn")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Brand(url *string, id *uint) {
	c := colly.NewCollector()
	data2 := []models.Brand{}
	brand2 := models.Brand{}
	c.OnHTML(".car-list", func(el *colly.HTMLElement) {
		el.ForEach(".car-list-item", func(_ int, ei *colly.HTMLElement) {
			// 车系名称
			name := ei.ChildText(".car-name")
			fmt.Printf("车系名称：%v\n", name)
			// 车系图片
			img := ei.ChildAttr(".relative > img", "src")
			if !strings.Contains(img, "http") {
				img = "http:" + img
			}
			fmt.Printf("车系名称：%v\n", img)
			// 品牌车型详情
			// href := "https:" + ei.ChildAttr("a", "href")
			// fmt.Println(href)
			brand2.Logo = img
			brand2.Name = name
			brand2.Pid = *id
			brand2.Type = TYPE
			data2 = append(data2, brand2)

			fmt.Printf("\n")
		})

		models.AddBrandBatch(&data2)
	})

	// 请求发起时回调,一般用来设置请求头等
	c.OnRequest(func(request *colly.Request) {
		currentProxyIndex := rand.Intn(52)
		request.ProxyURL = "http://" + logic.ProxyList[currentProxyIndex]
		fmt.Println("----> 开始请求了")
	})

	// 请求完成后回调
	c.OnResponse(func(response *colly.Response) {
		time.Sleep(time.Second * 2)
		fmt.Println("----> 开始返回了")
	})

	//请求发生错误回调
	c.OnError(func(response *colly.Response, err error) {
		fmt.Printf("发生错误了:%v", err)
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
	})

	// 在所有OnHTML之后执行，可以用来做一些回收操作
	c.OnScraped(func(response *colly.Response) {
		wait.Done()
		fmt.Println("----> 所有匹配已完成")
	})

	// 开始访问
	err := c.Visit(*url)
	if err != nil {
		log.Fatalln(err)
	}

	// return data2
}

// 车型菜单
func CarMenu() {
	c := colly.NewCollector()
	c.OnHTML(".car-model-list-warp", func(em *colly.HTMLElement) {
		em.ForEach("ul", func(_ int, eu *colly.HTMLElement) {
			eu.ForEach("li", func(_ int, el *colly.HTMLElement) {
				// 获取车型名称
				name := el.ChildText(".model-name")
				fmt.Println(name)

			})
			fmt.Printf("\n\n")
		})
	})

	// 请求发起时回调,一般用来设置请求头等
	c.OnRequest(func(request *colly.Request) {
		fmt.Println("----> 开始请求了")
	})

	// 请求完成后回调
	c.OnResponse(func(response *colly.Response) {
		fmt.Println("----> 开始返回了")
	})

	//请求发生错误回调
	c.OnError(func(response *colly.Response, err error) {
		fmt.Printf("发生错误了:%v", err)
	})

	// 在所有OnHTML之后执行，可以用来做一些回收操作
	c.OnScraped(func(response *colly.Response) {
		fmt.Println("----> 所有匹配已完成")
	})

	// 开始访问
	err := c.Visit("https://price.pcauto.com.cn/sg3996/#ad=21169")
	if err != nil {
		log.Fatalln(err)
		return
	}
}

// 车型菜单-模拟点击更多
func CarMenuMoreClick() {

}

// 数据持久化
// func csvSave(fName string, data []Article) error {
// 	file, err := os.Create(fName)
// 	if err != nil {
// 		log.Fatalf("Cannot create file %q: %s\n", fName, err)
// 	}
// 	defer file.Close()
// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()

// 	writer.Write([]string{"ID", "Title", "URL", "Created", "Reads", "Comments", "Feeds"})
// 	for _, v := range data {
// 		writer.Write([]string{strconv.Itoa(v.ID), v.Title, v.URL, v.Created, v.Reads, v.Comments, v.Feeds})
// 	}
// 	return nil
// }

// 重试机制
func retries(url *string, num int) {

}
