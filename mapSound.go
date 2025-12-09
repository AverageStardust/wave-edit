package main

import "math"

func mapSound(sample float64, location, width uint32) float64 {
	crunchedSample := sample

	for range 1 {
		crunchedSample = (math.Cos(crunchedSample) - 0.77) * 4.3
	}

	mixture := fadeInOut(float64(location), 100, float64(width))
	return fadeBetween(sample, crunchedSample, mixture)
}

func fadeBetween(x, y, fade float64) float64 {
	return x*(1-fade) + y*fade
}

func fadeInOut(x, edge float64, width float64) float64 {
	if x < edge {
		return fadeIn(x / edge)
	} else {
		return fadeIn((width - x) / edge)
	}
}

func fadeIn(x float64) float64 {
	if x < 0 {
		return 0
	} else if x > 1 {
		return 1
	} else {
		return 3*x*x - 2*x*x*x
	}
}
