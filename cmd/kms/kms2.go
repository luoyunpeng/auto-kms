package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var (
	lights1 = []time.Duration{
		time.Second * 0,                       //1,打开近光灯
		time.Second*7 + time.Microsecond*500,  // 2,夜间在没有照明或者照明条件不良的路上行驶
		time.Second*15 + time.Microsecond*500, // 3,通过拱桥
		time.Second * 21,                      // 4, 路口左转弯
		time.Second*27 + time.Microsecond*500, //5, 超车
		time.Second*35 + time.Microsecond*500, // 6,无路灯照明且对向车道150米内有车辆行驶
		time.Second*43 + time.Microsecond*800, // 7,夜间通过没有交通信号灯控制的路口
		time.Second*51 + time.Microsecond*500, // 8,在有路灯的道路上行驶
		time.Second*57 + time.Microsecond*500, // 9,通过急弯
		time.Second*63 + time.Microsecond*500, // 10,通过人行横道
		time.Second*69 + time.Microsecond*500, // 11,通过有信号灯指示的路口
		time.Second*76 + time.Microsecond*500, // 12,同方向近距离条件下你紧跟前车行驶
		time.Second*85 + time.Microsecond*500, // 13,雾天行驶
		time.Second*91 + time.Microsecond*500, // 14,路边临时停车
		time.Second * 98,                      // 15,大雨天行驶
		time.Second * 104,                     // 16,夜间模拟驾驶结束
	}

	audioLenMap1 = map[int]time.Duration{
		0:  time.Second * 4,
		1:  time.Second * 6,
		2:  time.Second * 3,
		3:  time.Second*3 + time.Millisecond*500,
		4:  time.Second * 2,
		5:  time.Second * 5,
		6:  time.Second * 6,
		7:  time.Second*4 - time.Millisecond*200,
		8:  time.Second * 3,
		9:  time.Second * 3,
		10: time.Second * 4,
		11: time.Second * 4,
		12: time.Second * 3,
		13: time.Second * 3,
		14: time.Second * 3,
		15: time.Second * 3,
	}
)

var (
	kmsAudio     beep.StreamSeekCloser
	kmsAudioInfo beep.Format
	playButton   *widget.Button
)

func main() {
	f, err := os.Open("../../asset/mp3/kms.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	kmsAudio = streamer
	kmsAudioInfo = format
	defer streamer.Close()
	defer speaker.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	initWindow()
}

func InitPlayButton1() *widget.Button {
	playButton = widget.NewButton("Order Play", nil)
	playButton.OnTapped = func() {
		playButton.Disable()
		play()
	}

	return playButton
}

func play1() {
	//stream.Seek(kmsAudioInfo.SampleRate.N(lights[0]))
	audios := randomAudioIndex(1, 14)
	log.Println("random:", audios)

	for _, ai := range audios {
		kmsAudio.Seek(kmsAudioInfo.SampleRate.N(lights[ai]))
		log.Println("play", ai)
		go P()
		time.Sleep(audioLenMap1[ai])
		if ai != 15 {
			speaker.Clear()
			time.Sleep(time.Second * 5)
		}
	}

}

func P1() {
	speaker.Play(beep.Seq(kmsAudio, beep.Callback(func() {
		playButton.Enable()
	})))
}

func initWindow1() {
	kmsExam := app.New()
	myWindow := kmsExam.NewWindow("科目三")

	kmsLabel := widget.NewLabelWithStyle("night light", 1, fyne.TextStyle{Bold: true})

	InitPlayButton()

	//	progress := widget.NewProgressBar()
	myWindow.SetContent(container.NewVBox(
		kmsLabel,
		//container.NewVBox(progress),
		playButton,
	))

	/*go func() {
		num := 0.0
		for num < 1.0 {
			time.Sleep(50 * time.Millisecond)
			progress.SetValue(num)
			num += 0.01
		}

		progress.SetValue(1)
	}()*/

	myWindow.Resize(fyne.NewSize(600, 800))
	myWindow.SetMaster()
	myWindow.ShowAndRun()
}

// random audio
func randomAudioIndex1(start int, end int) []int {
	// range check
	if end < start || (end-start) < 5 {
		return nil
	}

	nums := make([]int, 1, 7)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < 6 {
		//generate random num
		num := r.Intn(end-start) + start
		// check repeat
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}
		if !exist {
			nums = append(nums, num)
		}
	}
	nums = append(nums, 15)
	return nums
}
