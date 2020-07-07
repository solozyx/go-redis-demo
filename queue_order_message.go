package main

import (
	"encoding/json"
	"log"
	"time"
)

type OrderMessage struct {
	OrderID string `json:"order_id"`
}

func NewOrderMessage(orderId string) *OrderMessage {
	return &OrderMessage{OrderID: orderId}
}

func PublishOrderMessage(message *OrderMessage) error {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return rdsClient.LPush(_orderMsgQueue, string(msgBytes)).Err()
}

func ConsumeOrderMessage() error {
	strs, err := rdsClient.BRPop(10*time.Second, _orderMsgQueue).Result()
	// strs,err := rdsClient.BRPop(0,_orderMsgQueue).Result()
	if err != nil {
		log.Printf("ConsumeOrderMessage err=%v", err.Error())
	}
	if len(strs) > 0 {
		for i, str := range strs {
			log.Printf("str[%d]=%s\n", i, str)
		}
	}

	return nil
}
