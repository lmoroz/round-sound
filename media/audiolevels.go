package media

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
)

const (
	BandCount   = 64
	RefreshRate = 60 // FPS
)

// AudioLevelCapture manages audio level capture using WASAPI
type AudioLevelCapture struct {
	mu          sync.RWMutex
	isCapturing bool
	stopChan    chan struct{}
	callback    func([]float32)
	levels      []float32
	smoothing   float32
	rng         *rand.Rand
}

// NewAudioLevelCapture creates new audio capture instance
func NewAudioLevelCapture(callback func([]float32)) *AudioLevelCapture {
	return &AudioLevelCapture{
		callback:  callback,
		levels:    make([]float32, BandCount),
		smoothing: 0.3,
		stopChan:  make(chan struct{}),
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Start begins capturing audio levels
func (a *AudioLevelCapture) Start() error {
	a.mu.Lock()
	if a.isCapturing {
		a.mu.Unlock()
		return fmt.Errorf("audio capture already running")
	}
	a.isCapturing = true
	a.mu.Unlock()

	go a.captureLoop()
	log.Println("[AudioLevels] Capture started")
	return nil
}

// Stop stops capturing audio levels
func (a *AudioLevelCapture) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.isCapturing {
		return
	}

	close(a.stopChan)
	a.isCapturing = false
	log.Println("[AudioLevels] Capture stopped")
}

// captureLoop is the main capture loop
func (a *AudioLevelCapture) captureLoop() {
	// Initialize COM for this goroutine
	if err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		log.Printf("[AudioLevels] Failed to initialize COM: %v", err)
		return
	}
	defer ole.CoUninitialize()

	log.Println("[AudioLevels] COM initialized, starting capture loop")

	ticker := time.NewTicker(time.Second / RefreshRate)
	defer ticker.Stop()

	for {
		select {
		case <-a.stopChan:
			return
		case <-ticker.C:
			if err := a.captureFrame(); err != nil {
				// Silently continue on errors, try next frame
				continue
			}
		}
	}
}

// captureFrame captures audio meter reading
func (a *AudioLevelCapture) captureFrame() error {
	// Get default audio endpoint enumerator
	var mmde *wca.IMMDeviceEnumerator
	if err := wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return fmt.Errorf("failed to create device enumerator: %w", err)
	}
	defer mmde.Release()

	// Get default audio output device
	var mmd *wca.IMMDevice
	if err := mmde.GetDefaultAudioEndpoint(wca.ERender, wca.EConsole, &mmd); err != nil {
		return fmt.Errorf("failed to get default endpoint: %w", err)
	}
	defer mmd.Release()

	// Activate IAudioMeterInformation
	var objMeter *wca.IAudioMeterInformation
	if err := mmd.Activate(wca.IID_IAudioMeterInformation, wca.CLSCTX_ALL, nil, &objMeter); err != nil {
		return fmt.Errorf("failed to activate meter: %w", err)
	}
	defer objMeter.Release()

	// Get peak value
	var peak float32
	if err := objMeter.GetPeakValue(&peak); err != nil {
		return fmt.Errorf("failed to get peak: %w", err)
	}

	// Generate levels from peak
	a.generateLevelsFromPeak(peak)

	return nil
}

// generateLevelsFromPeak creates band levels from overall peak
func (a *AudioLevelCapture) generateLevelsFromPeak(peak float32) {
	levels := make([]float32, BandCount)

	// Apply gain boost for better visualization
	// Windows audio meter tends to give very low values
	const gainMultiplier = 10.0
	boostedPeak := peak * gainMultiplier

	// Clamp boosted peak to [0, 1]
	if boostedPeak > 1.0 {
		boostedPeak = 1.0
	}

	// Distribute peak across bands with variation for visual effect
	for i := 0; i < BandCount; i++ {
		// Create frequency-like distribution
		// Lower frequencies (bass) typically have more energy
		freq := float64(i) / float64(BandCount)

		// Bass boost (lower frequencies)
		bassBoost := float32(math.Exp(-freq * 2.0))

		// Random variation for organic feel
		variation := float32(a.rng.Float64()*0.3 - 0.15)

		// Combine
		levels[i] = boostedPeak * (0.5 + bassBoost*0.3 + variation)

		// Add some time-based animation
		timePhase := float32(math.Sin(float64(time.Now().UnixNano())/1e9*2.0 + float64(i)*0.1))
		levels[i] += timePhase * boostedPeak * 0.1

		// Add minimum level for visibility even when silent
		levels[i] += 0.05

		// Clamp to [0, 1]
		if levels[i] < 0 {
			levels[i] = 0
		}
		if levels[i] > 1 {
			levels[i] = 1
		}
	}

	a.updateLevels(levels)
}

// generateLevelsFromChannels creates band levels from channel peaks
func (a *AudioLevelCapture) generateLevelsFromChannels(channels []float32) {
	// Average channels
	var avgPeak float32
	for _, ch := range channels {
		avgPeak += ch
	}
	if len(channels) > 0 {
		avgPeak /= float32(len(channels))
	}

	a.generateLevelsFromPeak(avgPeak)
}

// updateLevels applies smoothing and sends levels to callback
func (a *AudioLevelCapture) updateLevels(newLevels []float32) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Apply exponential smoothing
	for i := 0; i < BandCount && i < len(newLevels); i++ {
		a.levels[i] = a.levels[i]*(1.0-a.smoothing) + newLevels[i]*a.smoothing
	}

	// Send to callback
	if a.callback != nil {
		levelsCopy := make([]float32, BandCount)
		copy(levelsCopy, a.levels)
		a.callback(levelsCopy)
	}
}
