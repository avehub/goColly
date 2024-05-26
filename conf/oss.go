package conf

var ossConf = map[string]map[string]string{
	// 需要读取配置文件
}

func GetAllOssConf() map[string]map[string]string {
	return ossConf
}

func GetOssConf(key string) map[string]string {
	if key == "" {
		key = "default"
	}
	return ossConf[key]
}
