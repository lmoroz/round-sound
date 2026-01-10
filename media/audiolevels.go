package media

import (
	"fmt"
	"log"
	"sync"
	"time"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
)

const (
	BandCount   = 64
	RefreshRate = 60
)

type AudioLevelCapture struct {
	mu          sync.RWMutex
	isCapturing bool
	stopChan    chan struct{}
	callback    func([]float32)
	config      FFTConfig
	buffer      []float32
	bufferSize  int
}

func NewAudioLevelCapture(callback func([]float32)) *AudioLevelCapture {
	return &AudioLevelCapture{
		callback:   callback,
		stopChan:   make(chan struct{}),
		config:     DefaultFFTConfig(),
		bufferSize: 0,
	}
}

func (a *AudioLevelCapture) UpdateConfig(fftSize int, freqMin, freqMax float64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.config.FFTSize = fftSize
	a.config.FreqMin = freqMin
	a.config.FreqMax = freqMax

	log.Printf("[AudioLevels] Config updated: FFTSize=%d, FreqMin=%.1f, FreqMax=%.1f", fftSize, freqMin, freqMax)
}

func (a *AudioLevelCapture) Start() error {
	a.mu.Lock()
	if a.isCapturing {
		a.mu.Unlock()
		return fmt.Errorf("audio capture already running")
	}
	a.isCapturing = true
	a.mu.Unlock()

	go a.captureLoop()
	log.Println("[AudioLevels] FFT Capture started")
	return nil
}

func (a *AudioLevelCapture) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.isCapturing {
		return
	}

	close(a.stopChan)
	a.isCapturing = false
	log.Println("[AudioLevels] FFT Capture stopped")
}

func (a *AudioLevelCapture) captureLoop() {
	if err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		log.Printf("[AudioLevels] Failed to initialize COM: %v", err)
		return
	}
	defer ole.CoUninitialize()

	log.Println("[AudioLevels] COM initialized, starting FFT capture loop")

	ticker := time.NewTicker(time.Second / RefreshRate)
	defer ticker.Stop()

	var mmde *wca.IMMDeviceEnumerator
	var mmd *wca.IMMDevice
	var audioClient *wca.IAudioClient
	var captureClient *wca.IAudioCaptureClient
	var pwfx *wca.WAVEFORMATEX
	var sampleRate uint32

	defer func() {
		if captureClient != nil {
			captureClient.Release()
		}
		if audioClient != nil {
			audioClient.Stop()
			audioClient.Release()
		}
		if mmd != nil {
			mmd.Release()
		}
		if mmde != nil {
			mmde.Release()
		}
		if pwfx != nil {
			ole.CoTaskMemFree(uintptr(unsafe.Pointer(pwfx)))
		}
	}()

	if err := wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		log.Printf("[AudioLevels] Failed to create device enumerator: %v", err)
		return
	}

	if err := mmde.GetDefaultAudioEndpoint(wca.ERender, wca.EConsole, &mmd); err != nil {
		log.Printf("[AudioLevels] Failed to get default endpoint: %v", err)
		return
	}

	if err := mmd.Activate(wca.IID_IAudioClient, wca.CLSCTX_ALL, nil, &audioClient); err != nil {
		log.Printf("[AudioLevels] Failed to activate audio client: %v", err)
		return
	}

	if err := audioClient.GetMixFormat(&pwfx); err != nil {
		log.Printf("[AudioLevels] Failed to get mix format: %v", err)
		return
	}

	sampleRate = pwfx.NSamplesPerSec
	log.Printf("[AudioLevels] Sample rate: %d Hz, Channels: %d, BitsPerSample: %d",
		sampleRate, pwfx.NChannels, pwfx.WBitsPerSample)

	hnsRequestedDuration := wca.REFERENCE_TIME(10000000)
	if err := audioClient.Initialize(
		wca.AUDCLNT_SHAREMODE_SHARED,
		wca.AUDCLNT_STREAMFLAGS_LOOPBACK,
		hnsRequestedDuration,
		0,
		pwfx,
		nil,
	); err != nil {
		log.Printf("[AudioLevels] Failed to initialize audio client with LOOPBACK: %v", err)
		return
	}

	if err := audioClient.GetService(wca.IID_IAudioCaptureClient, &captureClient); err != nil {
		log.Printf("[AudioLevels] Failed to get capture client: %v", err)
		return
	}

	if err := audioClient.Start(); err != nil {
		log.Printf("[AudioLevels] Failed to start audio client: %v", err)
		return
	}

	log.Println("[AudioLevels] WASAPI Loopback started successfully")

	for {
		select {
		case <-a.stopChan:
			return
		case <-ticker.C:
			a.processAudioFrame(captureClient, pwfx, sampleRate)
		}
	}
}

func (a *AudioLevelCapture) processAudioFrame(captureClient *wca.IAudioCaptureClient, pwfx *wca.WAVEFORMATEX, sampleRate uint32) {
	var packetLength uint32
	if err := captureClient.GetNextPacketSize(&packetLength); err != nil {
		return
	}

	if packetLength == 0 {
		a.sendSilence()
		return
	}

	var pData *byte
	var numFrames uint32
	var flags uint32

	for packetLength > 0 {
		if err := captureClient.GetBuffer(&pData, &numFrames, &flags, nil, nil); err != nil {
			return
		}

		if flags&wca.AUDCLNT_BUFFERFLAGS_SILENT != 0 || numFrames == 0 {
			captureClient.ReleaseBuffer(numFrames)
			a.sendSilence()
		} else {
			samples := a.extractSamples(pData, numFrames, pwfx)
			captureClient.ReleaseBuffer(numFrames)
			a.sendFFTLevels(samples, sampleRate)
		}

		if err := captureClient.GetNextPacketSize(&packetLength); err != nil {
			break
		}
	}
}

func (a *AudioLevelCapture) extractSamples(pData *byte, numFrames uint32, pwfx *wca.WAVEFORMATEX) []float32 {
	channels := pwfx.NChannels
	totalSamples := int(numFrames) * int(channels)

	samples := make([]float32, int(numFrames))

	if pwfx.WBitsPerSample == 32 {
		floatData := (*[1 << 30]float32)(unsafe.Pointer(pData))[:totalSamples:totalSamples]

		for i := 0; i < int(numFrames); i++ {
			var channelSum float32
			for ch := 0; ch < int(channels); ch++ {
				channelSum += floatData[i*int(channels)+ch]
			}
			samples[i] = channelSum / float32(channels)
		}
	} else if pwfx.WBitsPerSample == 16 {
		int16Data := (*[1 << 30]int16)(unsafe.Pointer(pData))[:totalSamples:totalSamples]

		for i := 0; i < int(numFrames); i++ {
			var channelSum int32
			for ch := 0; ch < int(channels); ch++ {
				channelSum += int32(int16Data[i*int(channels)+ch])
			}
			samples[i] = float32(channelSum) / float32(channels) / 32768.0
		}
	}

	return samples
}

func (a *AudioLevelCapture) sendFFTLevels(samples []float32, sampleRate uint32) {
	a.mu.RLock()
	config := a.config
	a.mu.RUnlock()

	a.buffer = append(a.buffer, samples...)

	if len(a.buffer) >= config.FFTSize {
		levels := ProcessFFT(a.buffer[:config.FFTSize], sampleRate, config, BandCount)

		a.buffer = a.buffer[len(samples):]
		if len(a.buffer) > config.FFTSize*2 {
			a.buffer = a.buffer[len(a.buffer)-config.FFTSize:]
		}

		if a.callback != nil {
			a.callback(levels)
		}
	}
}

func (a *AudioLevelCapture) sendSilence() {
	silence := make([]float32, BandCount)
	for i := range silence {
		silence[i] = 0.05
	}

	if a.callback != nil {
		a.callback(silence)
	}
}
