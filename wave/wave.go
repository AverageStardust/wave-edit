package wave

import (
	"errors"
	"io"
	"wave-edit/riff"
)

type WaveFile struct {
	Fmt  *FmtChunk
	Fact *FactChunk
	Data DataChunk
}

var ErrMissingFmt = errors.New("wave file missing format chunk")
var ErrMissingData = errors.New("wave file missing data chunk")

func init() {
	riff.RegisterRiffForm("WAVE", deserializeWave)
}

func CreateWave(format WaveFormat, channels uint16, samplesPerSec uint32) *WaveFile {
	factChunk := FactChunk(0)

	return &WaveFile{
		Fmt: &FmtChunk{
			Format:        format,
			Channels:      channels,
			SamplesPerSec: samplesPerSec,
		},
		Fact: &factChunk,
		Data: DataChunk{},
	}
}

func deserializeWave(reader io.Reader, size uint32) (riff.Chunk, error) {
	wave := &WaveFile{}

	for size > 0 {
		chunkSize, err := deserializeWaveChunk(reader, wave)
		if err != nil {
			return nil, err
		} else if chunkSize > size {
			return nil, riff.ErrReadTooMuch
		}

		size -= chunkSize
	}

	if wave.Fmt.Format == UNKNOWN_FORMAT {
		return nil, ErrMissingFmt
	} else if wave.Data == nil {
		return nil, ErrMissingData
	}

	if wave.Fact == nil {
		sampleCount := uint32(len(wave.Data)) / uint32(wave.Fmt.BlockSize())
		chunk := FactChunk(sampleCount)
		wave.Fact = &chunk
	}

	return wave, nil
}
func deserializeWaveChunk(reader io.Reader, waveFile *WaveFile) (uint32, error) {
	chunk, err := riff.DeserializeChunk(reader,
		func(reader io.Reader, id riff.FourCC, size uint32) (riff.Chunk, error) {
			switch id {
			case fmtChunkId:
				fmtChunk, err := fmtDeserializer(reader, id, size)
				waveFile.Fmt = fmtChunk
				return fmtChunk, err

			case factChunkId:
				factChunk, err := factDeserializer(reader, id, size)
				waveFile.Fact = factChunk
				return factChunk, err

			case dataChunkId:
				waveData, err := dataChunkDeserializer(reader, id, size)
				waveFile.Data = waveData
				return waveData, err

			default:
				return riff.IgnoreDeserializer(reader, id, size)
			}
		})

	if err != nil {
		return 0, err
	}

	return chunk.Size(), err
}

func (chunk *WaveFile) Serialize(writer io.Writer) error {
	err := riff.SerializeRiffHeader(writer, chunk.Size(), "WAVE")
	if err != nil {
		return err
	}

	err = chunk.Fmt.Serialize(writer)
	if err != nil {
		return err
	}

	if chunk.Fact != nil {
		err = chunk.Fact.Serialize(writer)
		if err != nil {
			return err
		}
	}

	err = chunk.Data.Serialize(writer)
	if err != nil {
		return err
	}

	return nil
}

func (chunk *WaveFile) Size() uint32 {
	var size uint32 = 12

	size += chunk.Fmt.Size()

	if chunk.Fact != nil {
		size += chunk.Fact.Size()
	}

	size += chunk.Data.Size()

	return size
}
