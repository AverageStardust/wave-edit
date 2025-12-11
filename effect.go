package main

import (
	"math"
	"wave-edit/wave"
)

func applyEffect(wave *wave.WaveFile, startTime, endTime, beatTime float64) error {
	start := uint32(startTime * float64(wave.Fmt.SamplesPerSec))
	end := uint32(endTime * float64(wave.Fmt.SamplesPerSec))
	beat := uint32(beatTime * float64(wave.Fmt.SamplesPerSec))

	for channel := range wave.Fmt.Channels {
		samples, err := wave.GetSamples(channel, start, end)
		if err != nil {
			return err
		}

		for n := range samples {
			// ignore last 2 samples
			if n < len(samples)-2 {
				samples[n] = effectSample(samples, uint32(n), end-start, beat)
			}
		}

		smoothSamples := make([]float64, end-start)
		for n := range samples {
			smoothSamples[n] = samples[n]

			// ignore first 2 and last two samples
			if n > 1 && n < len(samples)-2 {
				// detect clipping
				if math.Abs(samples[n-2]-samples[n+2]) > 0.3 {
					// smooth samples
					smoothSamples[n] = (samples[n-2] + samples[n-1] + samples[n] + samples[n+1] + samples[n+2]) / 5
				}
			}
		}

		err = wave.SetSamples(channel, start, smoothSamples)
		if err != nil {
			return err
		}
	}

	return nil
}

func effectSample(samples []float64, location, width, beat uint32) float64 {
	if location < width/2 {
		// for the first half of the effected region
		return speedUpEffect(samples, location, width)
	} else {
		// for the second half of the effected region
		return repeatEffect(samples, location, width, beat)
	}
}

func speedUpEffect(samples []float64, location, width uint32) float64 {
	// finds the index of sample 2x^2 + x
	// in other words, sound will play at an increasing speed
	x := float64(location) / float64(width)
	y := 2*(x*x) + x
	y *= float64(width)

	// smoothly mix between samples
	integer, fractional := math.Modf(y)
	leftSample := samples[uint32(integer)]
	rightSample := samples[uint32(integer)+1]
	sample := fadeBetween(leftSample, rightSample, fractional)

	return sample
}

func repeatEffect(samples []float64, location, width, beat uint32) float64 {
	// find the index where the last two beats start
	twoBeats := beat * 2
	lastTwoBeats := width - twoBeats

	// repeat the last two beats of sound over and over
	y := lastTwoBeats + location%twoBeats
	return samples[y]
}

func fadeBetween(x, y, fade float64) float64 {
	return x*(1-fade) + y*fade
}
