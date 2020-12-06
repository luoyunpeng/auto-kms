package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var (
	Light1  = time.Second * 0                       //1,打开近光灯
	Light2  = time.Second*7 + time.Microsecond*500  // 2,夜间在没有照明或者照明条件不良的路上行驶
	Light3  = time.Second*15 + time.Microsecond*500 // 3,通过拱桥
	Light4  = time.Second * 21                      // 4, 路口左转弯
	Light5  = time.Second*27 + time.Microsecond*500 //5, 超车
	Light6  = time.Second*35 + time.Microsecond*500 // 6,无路灯照明且对向车道150米内有车辆行驶
	Light7  = time.Second*43 + time.Microsecond*800 // 7,夜间通过没有交通信号灯控制的路口
	Light8  = time.Second*51 + time.Microsecond*500 // 8,在有路灯的道路上行驶
	Light9  = time.Second*57 + time.Microsecond*500 // 9,通过急弯
	Light10 = time.Second*63 + time.Microsecond*500 // 10,通过人行横道
	Light11 = time.Second*69 + time.Microsecond*500 // 11,通过有信号灯指示的路口
	Light12 = time.Second*76 + time.Microsecond*500 // 12,同方向近距离条件下你紧跟前车行驶
	Light13 = time.Second*85 + time.Microsecond*500 // 13,雾天行驶
	Light14 = time.Second*91 + time.Microsecond*500 // 14,路边临时停车
	Light15 = time.Second * 98                      // 15,大雨天行驶
	Light16 = time.Second * 104                     // 16,夜间模拟驾驶结束
)

func main() {
	f, err := os.Open("../asset/mp3/kms.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	ctx, cancel := context.WithCancel(context.Background())

	// loop := beep.Loop(3, streamer)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		log.Println("****** play kms audio done ********")
	})))

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
				speaker.Lock()
				fmt.Println(format.SampleRate.D(streamer.Position()).Round(time.Second), streamer.Position())
				speaker.Unlock()
			}
		}
	}()

	go func() {
		time.Sleep(time.Second * 1)
		streamer.Seek(format.SampleRate.N(Light16))
		time.Sleep(5 * time.Second)
		speaker.Clear()
	}()

	signalListen1()
	cancel()
	speaker.Close()
}

func signalListen1() {
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	fmt.Println("=====")
	fmt.Printf("reveive signal: %v\n", <-quit)
	fmt.Println("**** Graceful shutdown kms system ****")
}
