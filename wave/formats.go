package wave

import (
	"cmp"
	"encoding/binary"
	"math"
)

type WaveFormat int

const PCM_FORMAT_TAG = 0x0001
const IEEE_FLOAT_FORMAT_TAG = 0x0003

const (
	UNKNOWN_FORMAT WaveFormat = iota
	PCM_8
	PCM_16
	PCM_24
	PCM_32
	PCM_FLOAT32
	PCM_FLOAT64
)

func createWaveFormat(formatTag uint16, bitDepth uint16) WaveFormat {
	switch formatTag {
	case PCM_FORMAT_TAG:
		switch bitDepth {
		case 8: // unsigned 8-bit integer PCM
			return PCM_8
		case 16: // signed 16-bit integer PCM
			return PCM_16
		case 24: // signed 24-bit integer PCM
			return PCM_24
		case 32: // signed 32-bit integer PCM
			return PCM_32
		default:
			return UNKNOWN_FORMAT
		}

	case IEEE_FLOAT_FORMAT_TAG:
		switch bitDepth {
		case 32: // 32-bit IEEE Float PCM
			return PCM_FLOAT32
		case 64: // 64-bit IEEE Float PCM
			return PCM_FLOAT64
		default:
			return UNKNOWN_FORMAT
		}

	default:
		return UNKNOWN_FORMAT
	}
}

func (fmt WaveFormat) Properties() (formatTag uint16, byteDepth uint16) {
	switch fmt {
	case PCM_8:
		return PCM_FORMAT_TAG, 1
	case PCM_16:
		return PCM_FORMAT_TAG, 2
	case PCM_24:
		return PCM_FORMAT_TAG, 3
	case PCM_32:
		return PCM_FORMAT_TAG, 4
	case PCM_FLOAT32:
		return IEEE_FLOAT_FORMAT_TAG, 4
	case PCM_FLOAT64:
		return IEEE_FLOAT_FORMAT_TAG, 8
	default:
		panic("Unknown wave format")
	}
}

func (fmt WaveFormat) SampleGetter() func([]byte) float64 {
	switch fmt {
	case PCM_8:
		return getPCM8Sample
	case PCM_16:
		return getPCM16Sample
	case PCM_24:
		return getPCM24Sample
	case PCM_32:
		return getPCM32Sample
	case PCM_FLOAT32:
		return getPCMFloat32Sample
	case PCM_FLOAT64:
		return getPCMFloat64Sample
	default:
		panic("Unknown wave format")
	}
}

func (fmt WaveFormat) SampleSetter() func([]byte, float64) {
	switch fmt {
	case PCM_8:
		return setPCM8Sample
	case PCM_16:
		return setPCM16Sample
	case PCM_24:
		return setPCM24Sample
	case PCM_32:
		return setPCM32Sample
	case PCM_FLOAT32:
		return setPCMFloat32Sample
	case PCM_FLOAT64:
		return setPCMFloat64Sample
	default:
		panic("Unknown wave format")
	}
}

func getPCM8Sample(sampleData []byte) float64 {
	return (float64(sampleData[0]) - (1 << 7)) / (1 << 7)
}

func setPCM8Sample(sampleData []byte, sample float64) {
	clampedSample := clamp((sample+1)*(1<<7), 0, 255)
	sampleData[0] = uint8(clampedSample)
}

func getPCM16Sample(sampleData []byte) float64 {
	sampleInt := int16(binary.LittleEndian.Uint16(sampleData))
	return float64(sampleInt) / (1 << 15)
}

func setPCM16Sample(sampleData []byte, sample float64) {
	clampedSample := clamp(sample*(1<<15), -1<<15, 1<<15-1)
	binary.LittleEndian.PutUint16(sampleData, uint16(int16(clampedSample)))
}

func getPCM24Sample(sampleData []byte) float64 {
	var extendedSampleData [4]byte
	copy(extendedSampleData[1:4], sampleData)

	// sign extend
	if (extendedSampleData[1] & 0x80) > 0 {
		extendedSampleData[0] = 0xFF
	}

	sampleInt := int32(binary.LittleEndian.Uint32(extendedSampleData[:]))
	return float64(sampleInt) / (1 << 23)
}

func setPCM24Sample(sampleData []byte, sample float64) {
	clampedSample := clamp(sample*(1<<23), -1<<23, 1<<23-1)

	var extendedSampleData [4]byte
	binary.LittleEndian.PutUint32(extendedSampleData[:], uint32(int32(clampedSample)))

	copy(sampleData, extendedSampleData[1:4])
}

func getPCM32Sample(sampleData []byte) float64 {
	sampleInt := binary.LittleEndian.Uint32(sampleData)
	return float64(int32(sampleInt)) / (1 << 31)
}

func setPCM32Sample(sampleData []byte, sample float64) {
	clampedSample := clamp(sample*(1<<31), -1<<31, 1<<31-1)
	binary.LittleEndian.PutUint32(sampleData, uint32(int32(clampedSample)))
}

func getPCMFloat32Sample(sampleData []byte) float64 {
	sampleInt := binary.LittleEndian.Uint32(sampleData)
	return float64(math.Float32frombits(sampleInt))
}

func setPCMFloat32Sample(sampleData []byte, sample float64) {
	sampleBits := math.Float32bits(float32(sample))
	binary.LittleEndian.PutUint32(sampleData, sampleBits)
}

func getPCMFloat64Sample(sampleData []byte) float64 {
	sampleInt := binary.LittleEndian.Uint64(sampleData)
	return math.Float64frombits(sampleInt)
}

func setPCMFloat64Sample(sampleData []byte, sample float64) {
	sampleBits := math.Float64bits(sample)
	binary.LittleEndian.PutUint64(sampleData, sampleBits)
}

func clamp[T cmp.Ordered](value, minimum, maximum T) T {
	return max(min(value, maximum), minimum)
}
