package media

import (
	"math"
	"math/cmplx"

	"github.com/mjibson/go-dsp/fft"
)

type FFTConfig struct {
	FFTSize int
	FreqMin float64
	FreqMax float64
}

func DefaultFFTConfig() FFTConfig {
	return FFTConfig{
		FFTSize: 2048,
		FreqMin: 20,
		FreqMax: 20000,
	}
}

func ApplyHannWindow(samples []float64) {
	n := len(samples)
	for i := 0; i < n; i++ {
		samples[i] *= 0.5 * (1.0 - math.Cos(2.0*math.Pi*float64(i)/float64(n-1)))
	}
}

func ProcessFFT(samples []float32, sampleRate uint32, config FFTConfig, bandCount int) []float32 {
	if len(samples) < config.FFTSize {
		return make([]float32, bandCount)
	}

	float64Samples := make([]float64, config.FFTSize)
	for i := 0; i < config.FFTSize; i++ {
		float64Samples[i] = float64(samples[i])
	}

	ApplyHannWindow(float64Samples)

	fftResult := fft.FFTReal(float64Samples)

	halfSize := config.FFTSize / 2
	magnitudes := make([]float64, halfSize)
	for i := 0; i < halfSize; i++ {
		magnitudes[i] = cmplx.Abs(fftResult[i])
	}

	bands := groupIntoBands(magnitudes, sampleRate, config, bandCount)

	return bands
}

func groupIntoBands(magnitudes []float64, sampleRate uint32, config FFTConfig, bandCount int) []float32 {
	halfSize := len(magnitudes)
	freqPerBin := float64(sampleRate) / float64(len(magnitudes)*2)

	logMin := math.Log10(config.FreqMin)
	logMax := math.Log10(config.FreqMax)
	logStep := (logMax - logMin) / float64(bandCount)

	bands := make([]float32, bandCount)

	for band := 0; band < bandCount; band++ {
		freqStart := math.Pow(10, logMin+float64(band)*logStep)
		freqEnd := math.Pow(10, logMin+float64(band+1)*logStep)

		binStart := int(freqStart / freqPerBin)
		binEnd := int(freqEnd / freqPerBin)

		if binStart >= halfSize {
			binStart = halfSize - 1
		}
		if binEnd >= halfSize {
			binEnd = halfSize - 1
		}
		if binEnd <= binStart {
			binEnd = binStart + 1
		}

		var sum float64
		count := 0
		for bin := binStart; bin < binEnd; bin++ {
			if bin >= 0 && bin < halfSize {
				sum += magnitudes[bin]
				count++
			}
		}

		var avg float64
		if count > 0 {
			avg = sum / float64(count)
		}

		db := 20.0 * math.Log10(avg+1e-10)
		normalizedDB := (db + 60.0) / 60.0

		if normalizedDB < 0 {
			normalizedDB = 0
		}
		if normalizedDB > 1 {
			normalizedDB = 1
		}

		bands[band] = float32(normalizedDB)
	}

	return bands
}
