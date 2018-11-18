package main

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/beego/ms304w-client/basis/conf"
	l "github.com/beego/ms304w-client/basis/log"
	"github.com/beego/ms304w-client/controllers"
	_ "github.com/beego/ms304w-client/routers"
)

var log = l.New("main")

func main() {
	log.Info("start main...")
	beego.ErrorController(&controllers.ErrorController{})

	h := func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			log.Warn(origin)
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "X-HTTP-Method-Override, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// w.Header().Set("Access-Control-Max-Age", "86400")

		controllers.Server.ServeHTTP(w, r)
	}

	// socket
	http.HandleFunc("/socket.io/", h)

	// http.Handle("/socket.io/", controllers.Server)

	go func() {
		if err := http.ListenAndServe(conf.String("socket_url"), nil); err != nil {
			panic(err)
		}
	}()

	beego.SetStaticPath("/admin", "admin")
	beego.SetStaticPath("/client", "client")
	beego.SetStaticPath("/rfid", "rfid")

	beego.Run()

	log.Info("end main...")
}
