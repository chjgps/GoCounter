package syslog

import (
	"encoding/json"
	// "os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type config struct {
	Filename string `json:"filename"`
	MaxLines int    `json:"maxlines"`
	MaxSize  int    `json:"maxsize"`
	Daily    bool   `json:"daily"`
	MaxDays  int    `json:"maxdays"`
	Rotate   bool   `json:"rotate"`
	Level    int    `json:"level"`
	// Perm     os.FileMode `json:"perm"`
	Perm string `json:"perm"`
}

func newConfig(filename string) *config {
	return &config{
		Filename: filename,
		MaxLines: 1000000,
		MaxSize:  1 << 28,
		Daily:    true,
		MaxDays:  365,
		Rotate:   true,
		Level:    6,
		Perm:     "0777",
	}
}

func (c *config) String() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func New() *logs.BeeLogger {
	l := logs.NewLogger(10000)

	adapter := beego.AppConfig.String("log_adapter")

	conf := newConfig("logs/ms304w-client.log")
	if err := l.SetLogger(adapter, conf.String()); err != nil {
		panic(err)
	}

	l.SetLogFuncCallDepth(2)

	l.EnableFuncCallDepth(true)

	return l
}
