package aigenRecorder

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gen2brain/malgo"
	wave "github.com/zenwerk/go-wave"
)

// AdvancedVoiceRecorder provides enhanced recording with customizable duration
type AdvancedVoiceRecorder struct {
	SampleRate   int
	Channels     int
	Duration     time.Duration
	IsRecording  bool
	ctx          *malgo.AllocatedContext
	device       *malgo.Device
	samples      []byte
}

// NewAdvancedRecorder creates a new advanced recorder
func NewAdvancedRecorder() *AdvancedVoiceRecorder {
	return &AdvancedVoiceRecorder{
		SampleRate: 44100,
		Channels:   1,
		Duration:   20 * time.Second,
	}
}

// StartRecording begins audio capture with a specified duration
func (r *AdvancedVoiceRecorder) StartRecording(duration time.Duration) (string, error) {
	if duration == 0 {
		duration = r.Duration
	}

	r.Duration = duration

	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		fmt.Printf("Audio: %v\n", message)
	})
	if err != nil {
		return "", fmt.Errorf("failed to initialize context: %w", err)
	}
	r.ctx = ctx

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = uint32(r.Channels)
	deviceConfig.SampleRate = uint32(r.SampleRate)
	deviceConfig.Alsa.NoMMap = 1

	capturedSampleCount := uint32(0)
	r.samples = make([]byte, 0)

	sizeInBytes := uint32(malgo.SampleSizeInBytes(deviceConfig.Capture.Format))
	
	onRecvFrames := func(outputBuffer, inputBuffer []byte, framecount uint32) {
		sampleCount := framecount * deviceConfig.Capture.Channels * sizeInBytes
		r.samples = append(r.samples, inputBuffer...)
		capturedSampleCount += sampleCount
	}

	captureCallbacks := malgo.DeviceCallbacks{
		Data: onRecvFrames,
	}

	device, err := malgo.InitDevice(ctx.Context, deviceConfig, captureCallbacks)
	if err != nil {
		return "", fmt.Errorf("failed to initialize device: %w", err)
	}
	r.device = device

	if err := device.Start(); err != nil {
		return "", fmt.Errorf("failed to start device: %w", err)
	}

	r.IsRecording = true
	fmt.Printf("🎤 Recording for %v...\n", duration)

	// Create context with timeout
	recordingCtx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	// Handle interruption
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-recordingCtx.Done():
		fmt.Println("✓ Recording completed")
	case <-sigChan:
		fmt.Println("⚠ Recording interrupted")
	}

	r.StopRecording()

	// Save recording
	filePath := r.saveRecording()
	return filePath, nil
}

// StopRecording stops the current recording
func (r *AdvancedVoiceRecorder) StopRecording() {
	if r.device != nil {
		r.device.Stop()
		r.device.Uninit()
		r.device = nil
	}
	if r.ctx != nil {
		r.ctx.Uninit()
	}
	r.IsRecording = false
}

// saveRecording saves the captured audio to a WAV file
func (r *AdvancedVoiceRecorder) saveRecording() string {
	filePath := randomName()
	
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return ""
	}
	defer f.Close()

	param := wave.WriterParam{
		Out:           f,
		Channel:       r.Channels,
		SampleRate:    r.SampleRate,
		BitsPerSample: 16,
	}

	w, err := wave.NewWriter(param)
	if err != nil {
		fmt.Printf("Error creating wave writer: %v\n", err)
		return ""
	}
	defer w.Close()

	_, err = w.Write(r.samples)
	if err != nil {
		fmt.Printf("Error writing samples: %v\n", err)
		return ""
	}

	fmt.Printf("✓ Recording saved to %s\n", filePath)
	return filePath
}

// GetRecordingStats returns statistics about the current recording
func (r *AdvancedVoiceRecorder) GetRecordingStats() map[string]interface{} {
	return map[string]interface{}{
		"sampleRate":       r.SampleRate,
		"channels":         r.Channels,
		"duration":         r.Duration.String(),
		"sampleCount":      len(r.samples),
		"isRecording":      r.IsRecording,
		"estimatedSeconds": float64(len(r.samples)) / float64(r.SampleRate*r.Channels*2),
	}
}

