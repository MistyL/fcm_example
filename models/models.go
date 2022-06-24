package models

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type Paramters struct {
	UserId   string `json:"userId"`
	Topic    string `json:"topic"`
	Platform string `json:"platform"`
}

type RegistryUser struct {
	UserId   string `json:"userId"`
	Token    string `json:"token"`
	Platform string `json:"platform"`
}

type SubscribeTopic struct {
	Topic  string   `json:"topic"`
	UserId []string `json:"userId"`
}

type ListResp struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Users   []string `json:"users,omitempty"`
	Topics  []string `json:"topics,omitempty"`
}

type SendMessage struct {
	Title     string                 `json:"title"`
	Condition Conditions             `json:"conditions"`
	Message   map[string]interface{} `json:"message"`
}

type Conditions struct {
	UserId   []string `json:"userId,omitempty"`
	Topic    string   `json:"topic,omitempty"`
	Platform string   `json:"platform,omitempty"`
}

func GetDefaltMessage(userId, topic, platform string) []byte {
	str := `{"data":{"key1":"paramter1","key2":"paramter2"},"notification":{"title":"test_1","body":"test body 1"}}`
	msg := make(map[string]interface{})
	json.Unmarshal([]byte(str), &msg)
	res := &SendMessage{
		Message: msg,
	}
	if userId != "" {
		res.Condition.UserId = []string{userId}
	}
	if topic != "" {
		res.Condition.Topic = topic
	}
	if platform != "" {
		res.Condition.Platform = platform
	}
	res.Title = randStr(10)
	fmt.Printf("=========res , %v\n", res)
	rr, _ := json.Marshal(res)
	return rr
}

func randStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}
