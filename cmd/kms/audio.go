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
		speaker.Close()
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
		time.Sleep(time.Second * 10)
		streamer.Seek(format.SampleRate.N(time.Second * 91))
	}()

	signalListen1()
	cancel()

}

func signalListen1() {
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	fmt.Println("=====")
	fmt.Printf("reveive signal: %v\n", <-quit)
	fmt.Println("**** Graceful shutdown kms system ****")
}
