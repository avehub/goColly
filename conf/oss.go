package conf

var ossConf = map[string]map[string]string{
	"default": {
		// 配置信息

	},
	"prod": {
		// 配置信息

	},
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
