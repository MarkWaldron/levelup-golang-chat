package user

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
)

const (
	usersKey       = "users"
	userChannelFmt = "user:%s:channels"
	ChannelsKey    = "channels"
)

type User struct {
	name            string
	channelsHandler *redis.PubSub
	stopListenerChan chan struct{}
	listening        bool
	MessageChan chan redis.Message
}