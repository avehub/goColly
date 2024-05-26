package common

import (
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var brucket string

type client interface {
}

func init() {
	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 设置连接数为10，每个主机的最大闲置连接数为20，每个主机的最大连接数为20。
	conn := oss.MaxConns(10, 20, 20)
	// 设置HTTP连接超时时间为20秒，HTTP读取或写入超时时间为60秒。
	time := oss.Timeout(20, 60)
	// 设置是否支持将自定义域名作为Endpoint，默认不支持。
	cname := oss.UseCname(true)
	// 设置HTTP的User-Agent头，默认为aliyun-sdk-go。
	userAgent := oss.UserAgent("aliyun-sdk-go")
	// 设置是否开启HTTP重定向，默认开启。
	redirect := oss.RedirectEnabled(true)
	// 设置是否开启SSL证书校验，默认关闭。
	verifySsl := oss.InsecureSkipVerify(false)
	// 设置代理服务器地址和端口。
	//proxy := oss.Proxy("yourProxyHost")
	// 设置代理服务器的主机地址和端口，代理服务器验证的用户名和密码。
	authProxy := oss.AuthProxy("yourProxyHost", "yourProxyUserName", "yourProxyPassword")
	// 开启CRC加密。
	crc := oss.EnableCRC(true)
	// 设置日志模式。
	logLevel := oss.SetLogLevel(oss.LogOff)

	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	// yourRegion填写Endpoint对应的Region信息，例如cn-hangzhou。

	client, err := oss.New("yourEndpoint", "", "", oss.SetCredentialsProvider(&provider), oss.AuthVersion(oss.AuthV4), oss.Region("yourRegion"))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	fmt.Printf("client:%#v\n", client)
}

func setBucket(bucket string) {

}
