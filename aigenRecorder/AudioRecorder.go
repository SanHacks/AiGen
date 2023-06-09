package aigenRecorder

import (
	"fmt"
	"github.com/gen2brain/malgo"
	wave "github.com/zenwerk/go-wave"
	"os"
	"time"
)

// VoiceRecorder takes voice/speech as input and uses OpenAI Whisper
// To Listen To Words Said In Captured Speech
// Not Sure If This Can Scale But It as Good as it is supposed to be
// For PERSONAL USAGE
// PLEASE DO NOT TOUCH THIS! I HAD TO PERFORM MIRACLES TO MAKE THIS WORK.
func VoiceRecorder() (string, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		//fmt.Printf("LOG <%v>\n", message)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = 44100
	deviceConfig.Alsa.NoMMap = 1

	var capturedSampleCount uint32
	pCapturedSamples := make([]byte, 0)

	sizeInBytes := uint32(malgo.SampleSizeInBytes(deviceConfig.Capture.Format))
	onRecvFrames := func(pSample2, pSample []byte, framecount uint32) {
		sampleCount := framecount * deviceConfig.Capture.Channels * sizeInBytes
		newCapturedSampleCount := capturedSampleCount + sampleCount
		pCapturedSamples = append(pCapturedSamples, pSample...)
		fmt.Println(capturedSampleCount, "/", newCapturedSampleCount, "samples captured.")
		capturedSampleCount = newCapturedSampleCount
	}

	fmt.Println("Recording for 10 seconds...")
	captureCallbacks := malgo.DeviceCallbacks{
		Data: onRecvFrames,
	}
	device, err := malgo.InitDevice(ctx.Context, deviceConfig, captureCallbacks)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = device.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Record for 10 seconds
	time.Sleep(10 * time.Second)

	device.Uninit()
	filePathName := randomName()
	f, err := os.Create(filePathName)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	param := wave.WriterParam{
		Out:           f,
		Channel:       1,
		SampleRate:    44100,
		BitsPerSample: 16,
	}
	w, err := wave.NewWriter(param)

	_, err = w.Write(pCapturedSamples)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Recording saved to", filePathName)

	defer w.Close()

	return filePathName, nil
}
