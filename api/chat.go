package api

import (
	"chat/user"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader websocket.Upgrader

var connectedUser = make(map[string]*user.User)

func H(rdb *redis.Client, fn func(http.ResponseWriter, *http.Request, *redis.Client)) func(httpResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, rdb)
	}
}



func ChatWebSocketHandler(w http.ResponseWriter, r *http.Request, rdb *redis.Client) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		handleWSError(err, conn)
		return
	}

	err = onConnect(r, conn, rdb)
	if err != nil {
		handleWSError(err, conn)
		return
	}

	closeCh := onDisconnect(r, conn, rdb)

	// on receive message from redis channels
	onChannelMessage(conn, r)

	loop:
		for {
			select {
			case <-closeCh:
				break loop
			default:
				onUserMessage(conn, r, rdb)
			}
		}
}