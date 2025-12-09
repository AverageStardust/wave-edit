package wave

import (
	"errors"
)

type SampleMapper func(in float64, location uint32, width uint32) (out float64)

var ErrChannelDoesNotExist = errors.New("accessed channel that does not exist")
var ErrSampleOutOfRange = errors.New("sample location not in file")
var ErrInvalidSampleRange = errors.New("sample range end < start")

func (wave *WaveFile) MapAllSamples(mapper SampleMapper) error {
	sampleCount := wave.Fact.Samples()

	for i := range wave.Fmt.Channels {
		err := wave.MapSamples(i, 0, sampleCount, mapper)
		if err != nil {
			return err
		}
	}

	return nil
}

func (wave *WaveFile) MapSamples(channel uint16, start uint32, end uint32, mapper SampleMapper) error {
	samples, err := wave.GetSamples(channel, start, end)
	if err != nil {
		return err
	}

	for n, sample := range samples {
		samples[n] = mapper(sample, uint32(n), end-start)
	}

	return wave.SetSamples(channel, start, samples)
}

func (wave *WaveFile) GetSample(channel uint16, location uint32) (float64, error) {
	samples, err := wave.GetSamples(channel, location, 1)
	if err != nil {
		return 0, err
	}

	return samples[0], nil
}

func (wave *WaveFile) GetSamples(channel uint16, start, end uint32) ([]float64, error) {
	if end > wave.Fact.Samples() {
		return nil, ErrSampleOutOfRange
	} else if end < start {
		return nil, ErrInvalidSampleRange
	} else if channel >= wave.Fmt.Channels {
		return nil, ErrChannelDoesNotExist
	}

	_, byteDepth := wave.Fmt.Format.Properties()
	blockSize := uint32(wave.Fmt.BlockSize())
	index := blockSize*start + uint32(byteDepth*channel)
	getter := wave.Fmt.Format.SampleGetter()

	length := end - start
	samples := make([]float64, length)
	for n := range length {
		sampleData := wave.Data[index : index+uint32(byteDepth)]
		index += blockSize

		samples[n] = getter(sampleData)
	}

	return samples, nil
}

func (wave *WaveFile) SetSample(channel uint16, location uint32, sample float64) error {
	return wave.SetSamples(channel, location, []float64{sample})
}

func (wave *WaveFile) SetSamples(channel uint16, location uint32, samples []float64) error {
	if location+uint32(len(samples)) > wave.Fact.Samples() {
		return ErrSampleOutOfRange
	} else if channel >= wave.Fmt.Channels {
		return ErrChannelDoesNotExist
	}

	_, byteDepth := wave.Fmt.Format.Properties()
	blockSize := uint32(wave.Fmt.BlockSize())
	index := blockSize*location + uint32(byteDepth*channel)
	setter := wave.Fmt.Format.SampleSetter()

	for _, sample := range samples {
		sampleData := wave.Data[index : index+uint32(byteDepth)]
		index += blockSize

		setter(sampleData, sample)
	}

	return nil
}
