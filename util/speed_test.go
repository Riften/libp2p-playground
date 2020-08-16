package util

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestSpeed(t *testing.T) {
	newctx, _ := context.WithTimeout(context.Background(), 15 * time.Second)
	recorder := NewSpeedRecorder(newctx, 10)
	rander := rand.New(rand.NewSource(time.Now().UnixNano()))
	go func() {
		for {
			recorder.AddStamp(time.Now(),rander.Intn(10000))
			time.Sleep(200*time.Nanosecond)
		}
	}()

	go func() {
		recorder.Start()
		fmt.Println("\nrecorder end")
	}()

	go func() {
		for {
			//fmt.Println("aaa")
			fmt.Println(fmt.Sprintf("瞬时速度：%.3f，\t平均速度：%.3f\r", recorder.InstantV(), recorder.AverageV()))
			//fmt.Printf("bbbbb")
			time.Sleep(time.Second)
		}
	}()
	time.Sleep(17*time.Second)
}
