package conf

import (
	"github.com/astaxie/beego"
)

func String(key string) string {
	return beego.AppConfig.String(key)
}

func Bool(key string) bool {
	b, err := beego.AppConfig.Bool(key)
	if err != nil {
		panic(err.Error() + key)
	}

	return b
}

func Int(key string) int {
	i, err := beego.AppConfig.Int(key)
	if err != nil {
		panic(err.Error() + key)
	}

	return i
}
