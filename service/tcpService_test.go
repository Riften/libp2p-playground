package service

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewTCPService(t *testing.T) {
	tcpService := NewTCPService(context.Background())
	err := tcpService.StartListen(10000)
	if err != nil {
		t.Fatal(err)
	}
	err =  tcpService.StartSend("127.0.0.1", 10000)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(50* time.Millisecond)
	senderId := tcpService.SendList()[0]
	receiverId := tcpService.ReceiveList()[0]
	sender, _ := tcpService.GetSender(senderId)
	receiver, _ := tcpService.GetReceiver(receiverId)

	go func() {
		for {
			fmt.Println(fmt.Sprintf("发送瞬时：%.3f，\t发送平均%.3f\t接收瞬时：%.3f\t接收平均%.3f (b/ms)",
				sender.Recorder.InstantV(), sender.Recorder.AverageV(),
				receiver.Recorder.InstantV(), receiver.Recorder.AverageV()))
			time.Sleep(time.Second)
		}
	}()

	time.Sleep(10 * time.Second)
}
