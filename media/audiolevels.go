package media

import (
	"errors"
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

	// AUDCLNT_E_DEVICE_INVALIDATED full HRESULT: severity(0x8) + facility AUDCLNT(0x889) + code(0x004).
	hrAudClntDeviceInvalidated uintptr = 0x88890004

	// reinitBackoffTicks throttles WASAPI re-initialization attempts to ~500ms when the
	// session is down, so we don't hammer COM at the full 60Hz refresh rate.
	reinitBackoffTicks = RefreshRate / 2

	// defaultDeviceCheckInterval polls the current default render endpoint every
	// ~166ms to detect Windows "default output" switches. WASAPI does NOT signal
	// AUDCLNT_E_DEVICE_INVALIDATED in that scenario when the old device is still
	// physically connected — the loopback session just goes silent. Polling the
	// default endpoint ID is the only reliable way to react.
	defaultDeviceCheckInterval = 10
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

// wasapiSession bundles the WASAPI stack needed for one loopback capture.
// Tied to a single render endpoint — when the user changes Windows default output,
// the whole session must be released and reopened.
//
// The device enumerator (mmde) is intentionally kept outside this struct: it has
// a longer lifecycle than a single capture session and is reused across re-inits
// to poll the current default endpoint ID.
type wasapiSession struct {
	mmd           *wca.IMMDevice
	audioClient   *wca.IAudioClient
	captureClient *wca.IAudioCaptureClient
	pwfx          *wca.WAVEFORMATEX
	sampleRate    uint32
	deviceID      string
}

func (s *wasapiSession) release() {
	if s == nil {
		return
	}
	if s.captureClient != nil {
		s.captureClient.Release()
		s.captureClient = nil
	}
	if s.audioClient != nil {
		s.audioClient.Stop()
		s.audioClient.Release()
		s.audioClient = nil
	}
	if s.mmd != nil {
		s.mmd.Release()
		s.mmd = nil
	}
	if s.pwfx != nil {
		ole.CoTaskMemFree(uintptr(unsafe.Pointer(s.pwfx)))
		s.pwfx = nil
	}
}

func openWASAPISession(mmde *wca.IMMDeviceEnumerator) (*wasapiSession, error) {
	s := &wasapiSession{}

	if err := mmde.GetDefaultAudioEndpoint(wca.ERender, wca.EConsole, &s.mmd); err != nil {
		s.release()
		return nil, fmt.Errorf("get default endpoint: %w", err)
	}

	if err := s.mmd.GetId(&s.deviceID); err != nil {
		s.release()
		return nil, fmt.Errorf("get device id: %w", err)
	}

	if err := s.mmd.Activate(wca.IID_IAudioClient, wca.CLSCTX_ALL, nil, &s.audioClient); err != nil {
		s.release()
		return nil, fmt.Errorf("activate audio client: %w", err)
	}

	if err := s.audioClient.GetMixFormat(&s.pwfx); err != nil {
		s.release()
		return nil, fmt.Errorf("get mix format: %w", err)
	}

	s.sampleRate = s.pwfx.NSamplesPerSec

	hnsRequestedDuration := wca.REFERENCE_TIME(10000000)
	if err := s.audioClient.Initialize(
		wca.AUDCLNT_SHAREMODE_SHARED,
		wca.AUDCLNT_STREAMFLAGS_LOOPBACK,
		hnsRequestedDuration,
		0,
		s.pwfx,
		nil,
	); err != nil {
		s.release()
		return nil, fmt.Errorf("initialize audio client (LOOPBACK): %w", err)
	}

	if err := s.audioClient.GetService(wca.IID_IAudioCaptureClient, &s.captureClient); err != nil {
		s.release()
		return nil, fmt.Errorf("get capture client: %w", err)
	}

	if err := s.audioClient.Start(); err != nil {
		s.release()
		return nil, fmt.Errorf("start audio client: %w", err)
	}

	log.Printf("[AudioLevels] WASAPI loopback opened: deviceID=%s, sampleRate=%d Hz, channels=%d, bitsPerSample=%d",
		s.deviceID, s.sampleRate, s.pwfx.NChannels, s.pwfx.WBitsPerSample)

	return s, nil
}

// currentDefaultDeviceID returns the ID of the current default render endpoint,
// without keeping any references — caller decides whether to act on it.
func currentDefaultDeviceID(mmde *wca.IMMDeviceEnumerator) (string, error) {
	var dev *wca.IMMDevice
	if err := mmde.GetDefaultAudioEndpoint(wca.ERender, wca.EConsole, &dev); err != nil {
		return "", err
	}
	defer dev.Release()

	var id string
	if err := dev.GetId(&id); err != nil {
		return "", err
	}
	return id, nil
}

// isDeviceInvalidated reports whether err signals that the bound audio endpoint
// is gone (device unplugged, default output switched in Windows, etc).
func isDeviceInvalidated(err error) bool {
	var oerr *ole.OleError
	if errors.As(err, &oerr) {
		return oerr.Code() == hrAudClntDeviceInvalidated
	}
	return false
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
	if err := wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		log.Printf("[AudioLevels] Failed to create device enumerator: %v", err)
		return
	}
	defer mmde.Release()

	var session *wasapiSession
	defer func() { session.release() }()

	session, err := openWASAPISession(mmde)
	if err != nil {
		log.Printf("[AudioLevels] Initial WASAPI open failed: %v", err)
	}

	ticksSinceReinitAttempt := 0
	ticksSinceDeviceCheck := 0

	for {
		select {
		case <-a.stopChan:
			return
		case <-ticker.C:
			if session == nil {
				ticksSinceReinitAttempt++
				if ticksSinceReinitAttempt >= reinitBackoffTicks {
					ticksSinceReinitAttempt = 0
					if s, err := openWASAPISession(mmde); err == nil {
						session = s
						a.buffer = a.buffer[:0]
						log.Println("[AudioLevels] WASAPI session re-opened")
					} else {
						log.Printf("[AudioLevels] Reinit attempt failed: %v", err)
					}
				}
				a.sendSilence()
				continue
			}

			ticksSinceDeviceCheck++
			if ticksSinceDeviceCheck >= defaultDeviceCheckInterval {
				ticksSinceDeviceCheck = 0
				if currentID, err := currentDefaultDeviceID(mmde); err == nil && currentID != session.deviceID {
					log.Printf("[AudioLevels] Default render endpoint changed (%s -> %s) — reopening WASAPI session", session.deviceID, currentID)
					session.release()
					session = nil
					ticksSinceReinitAttempt = 0
					a.sendSilence()
					continue
				}
			}

			if err := a.processAudioFrame(session); err != nil {
				if isDeviceInvalidated(err) {
					log.Println("[AudioLevels] Audio endpoint invalidated — reopening WASAPI session")
				} else {
					log.Printf("[AudioLevels] Capture error, reopening session: %v", err)
				}
				session.release()
				session = nil
				ticksSinceReinitAttempt = 0
				a.sendSilence()
			}
		}
	}
}

func (a *AudioLevelCapture) processAudioFrame(s *wasapiSession) error {
	var packetLength uint32
	if err := s.captureClient.GetNextPacketSize(&packetLength); err != nil {
		return err
	}

	if packetLength == 0 {
		a.sendSilence()
		return nil
	}

	var pData *byte
	var numFrames uint32
	var flags uint32

	for packetLength > 0 {
		if err := s.captureClient.GetBuffer(&pData, &numFrames, &flags, nil, nil); err != nil {
			return err
		}

		if flags&wca.AUDCLNT_BUFFERFLAGS_SILENT != 0 || numFrames == 0 {
			if err := s.captureClient.ReleaseBuffer(numFrames); err != nil {
				return err
			}
			a.sendSilence()
		} else {
			samples := a.extractSamples(pData, numFrames, s.pwfx)
			if err := s.captureClient.ReleaseBuffer(numFrames); err != nil {
				return err
			}
			a.sendFFTLevels(samples, s.sampleRate)
		}

		if err := s.captureClient.GetNextPacketSize(&packetLength); err != nil {
			return err
		}
	}
	return nil
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
