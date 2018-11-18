package controllers

import (
	"github.com/beego/ms304w-client/basis/errors"
	"github.com/googollee/go-socket.io"
)

var Server *socketio.Server

func init() {
	if Server != nil {
		return
	}

	server, err := socketio.NewServer(nil)
	if err != nil {
		panic(errors.As(err))
	}

	Server = server

	server.On("connection", func(so socketio.Socket) {
		log.Info("on connection: %s", so.Id())

		so.Join("login")

		/*
			so.On("message", func(msg string){
				log.Info("on connection news: %s", msg)

				// so.Emit("message", msg)

				server.BroadcastTo("login", "news", msg)
			})
		*/

		// so.Emit("connected", "connected!")
	})

	server.On("disconnection", func(so socketio.Socket) {
		log.Info("on disconnection: %s", so.Id())

		so.Leave("login")
	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Info("on error")
	})
}
