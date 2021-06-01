package main

import (
	"fmt"
	"github.peaut.limit/mq"
)

type TestPro struct {
	msgContent string
}

//实现发送者
func (t *TestPro) MsgContent() string {
	return t.msgContent
}

//实现接受者
func (t *TestPro) Consumer(dateByte []byte) error {
	fmt.Println(string(dateByte))
	return nil
}

func main() {
	msg := fmt.Sprintf("这是测试任务")
	t := &TestPro{
		msgContent: msg,
	}
	queueExchange := &mq.QueueExchange{
		QuName: "test.rabbit",
		RtKey:  "rabbit.key",
		ExName: "test.rabbit.mq",
		ExType: "direct",
	}
	mqp := mq.New(queueExchange)
	mqp.RegisterProducer(t)
	mqp.RegisterReceiver(t)
	mqp.Start()
}
