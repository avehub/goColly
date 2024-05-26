package moto58

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"spider/models"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
)

func init() {
	rand.NewSource(time.Now().UnixNano())
	log.SetPrefix("[moto58]")
}

// 58moto品牌列表Domain接口参数
type Param struct {
	Platform        string `json:"platform"`
	Version         string `json:"version"`
	BrandEnergyType string `json:"brandEnergyType"`
	BrandVersion    string `json:"brandVersion"`
}

// 58moto车型列表LIST接口参数
// 245?platform=11&version=3.59.0&provinceName=陕西省&cityName=咸阳市&onSale=1&page=1&rows=20&brandEnergyType=3&goodMinPrice&goodMaxPrice&goodType&goodMinVolume&goodMaxVolume&sortType=0
type Param2 struct {
	Platform string `json:"platform"`
	Version  string `json:"version"`
	// ProvinceName    string `json:"provinceName"`
	// CityName        string `json:"cityName"`
	OnSale          string `json:"onSale"`
	Page            string `json:"page"`
	Rows            string `json:"rows"`
	BrandEnergyType string `json:"brandEnergyType"`
	SortType        string `json:"sortType"`
}

// 58moto品牌列表需要提取字段
type Brand struct {
	aleph           string
	brandId         float64
	brandLogo       string
	brandName       string
	keywords        string
	spelling        string
	brandEnergyType float64
	goodId          float64
	goodName        string
	goodPic         string
	originGoodPic   string
	seriesId        float64
	seriesName      string
	CreateTime      int
	UpdateTime      int
	DeleteTime      int
}

// 58moto车型列表需要提取字段
type Good struct {
	brandId       string
	goodId        string
	goodName      string
	goodPic       string
	originGoodPic string
	seriesId      string
	seriesName    string
}

const (
	// 摩托车网站地址
	DOMAIN = "https://m.58moto.com/clientApi/carport/brand/v2/all/list"
	// 摩托车品牌详情列表
	LIST = "https://m.58moto.com/clientApi/carport/goods/v4/brand/"
	// 类型，摩托 = 2
	TYPE = 2
	// 设置最大重试次数
	MAX = 3
)

var wait sync.WaitGroup
var response map[string]interface{}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func GetBody(url *string) map[string]interface{} {
	log.Println("请求URL: ", *url)
	client := &http.Client{}
	res, err := client.Get(*url)
	log.Println("返回结果: ", res)
	checkErr(err)
	defer res.Body.Close()
	// 读取响应体
	body, err := io.ReadAll(res.Body)
	log.Println("结果解析: ", body)
	checkErr(err)
	err = json.Unmarshal(body, &response)
	checkErr(err)

	return response
}

func getParameter(obj interface{}) string {
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)
	var t, str string
	// num := int(unsafe.Sizeof(param))
	num := typ.NumField()
	for i := 0; i < num; i++ {
		if i == 0 {
			t = "?"
		} else {
			t = "&"
		}

		field := typ.Field(i)
		jsonField := field.Tag.Get("json")
		str += fmt.Sprintf("%s%s=%v", t, jsonField, val.Field(i))
	}

	return str
}

func Start() {
	url := DOMAIN + getParameter(Param{"11", "3.59.0", "1", "1598596364"})
	log.Println("Start请求开始")
	body := GetBody(&url)
	log.Println("Start请求结束")
	data := body["data"].([]interface{})
	// resBrand := []map[string]interface{}{}
	resBrand := []models.Motos{}
	keys := reflect.TypeOf(Brand{})
	var status models.Motos
	for j, item := range data {
		// 因为服务器http请求并发问题，每20次请求限制睡眠5秒
		if j%20 == 0 {
			time.Sleep(5 * time.Second)
		}

		if itemMap, ok := item.(map[string]interface{}); ok {
			// log.Printf("%v\n", itemMap)
			brand := map[string]interface{}{}
			for i := 0; i < keys.NumField(); i++ {
				field := keys.Field(i)
				key := field.Name
				brand[key] = itemMap[key]
			}
			log.Printf("采集数据brand: %v\n", brand)

			// 使用mapstructure库进行转换
			err := mapstructure.Decode(brand, &status)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			if brandId, ok := brand["brandId"].(float64); ok {
				if brandEnergyType, ok := brand["brandEnergyType"].(float64); ok {
					wait.Add(1)
					go GoodList(brandId, brandEnergyType)
				}
			}
			resBrand = append(resBrand, status)

		}
	}

	log.Printf("%v", resBrand)
	models.AddMotoBatch(&resBrand)
	wait.Wait()
}

func GoodList(brandId float64, brandEnergyType float64) {
	defer wait.Done()

	time.Sleep(time.Second)
	url := LIST + fmt.Sprintf("%.0f", brandId) + getParameter(Param2{"11", "3.59.0", "1", "1", "1000", fmt.Sprintf("%.0f", brandEnergyType), "0"})
	log.Println("GoodList请求开始")
	body := GetBody(&url)
	log.Println("GoodList请求结束")
	resGood := []models.Motos{}
	data, ok := body["data"].([]interface{})
	if !ok {
		log.Println("GoodList中data数据异常")
		log.Println("body : ", body)
		return
	}
	// data := body["data"].([]interface{})
	log.Printf("GoodList数据解析%v", data)
	// resGood := []map[string]interface{}{}
	keys := reflect.TypeOf(Brand{})
	var status models.Motos
	for _, item := range data {
		if itemMap, ok := item.(map[string]interface{}); ok {
			log.Printf("%v\n", itemMap)
			good := map[string]interface{}{}
			for i := 0; i < keys.NumField(); i++ {
				field := keys.Field(i)
				key := field.Name
				good[key] = itemMap[key]
			}
			log.Printf("采集数据good: %v\n", good)
			// 使用mapstructure库进行转换
			err := mapstructure.Decode(good, &status)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			resGood = append(resGood, status)
		}

	}
	models.AddMotoBatch(&resGood)

	log.Printf("采集数据goods: %v\n", resGood)
}
