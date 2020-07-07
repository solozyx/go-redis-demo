package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	_orderMsgQueue = "balance:order:queue"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if rdsClient == nil {
		rdsClient = NewClient()
	}
	pong, err := rdsClient.Ping().Result()
	fmt.Println(pong, err)
}

func main() {
	// GetUserRemain(2)
	ReloadUserRemain()
}

func queue() {
	msg1 := NewOrderMessage("order_1")
	_ = PublishOrderMessage(msg1)
	for {
		_ = ConsumeOrderMessage()
	}
}
