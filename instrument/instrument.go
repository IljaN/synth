package instrument

import (
	"github.com/IljaN/synth/generator"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"os"
)

type Instrument struct {
	SamplingRate int
	BitRate      int
	Duration     int
}

type Filter interface {
	Filter(buf *audio.FloatBuffer)
}

func New(samplingRate int, bitRate int, duration int) *Instrument {
	p := &Instrument{
		SamplingRate: samplingRate,
		BitRate:      bitRate,
		Duration:     duration,
	}
	return p
}

func (in *Instrument) Out(o Out) *audio.FloatBuffer {
	buf := in.EmptyBuf()
	o(buf)
	return buf
}

func (in *Instrument) EmptyBuf() *audio.FloatBuffer {
	return &audio.FloatBuffer{
		Data:   make([]float64, in.Duration*2),
		Format: audio.FormatStereo48000,
	}
}

type Out func(buf *audio.FloatBuffer)

func OscOutput(osc *generator.Osc) func(buf *audio.FloatBuffer) {
	o := osc
	return func(buf *audio.FloatBuffer) {
		o.Fill(buf)
	}
}

func Channel(stages ...Out) Out {
	return func(buf *audio.FloatBuffer) {
		for i := range stages {
			stages[i](buf)
		}
	}
}

func Div(val float64) Out {
	return func(buf *audio.FloatBuffer) {
		for i := range buf.Data {
			if buf.Data[i] != 0.0 {
				buf.Data[i] /= val
			}
		}
	}
}

func (in *Instrument) Mix(inputs ...Out) Out {
	oscCount := float64(len(inputs))
	oscOuts := make([]*audio.FloatBuffer, 0, int(oscCount))

	// Pre allocate out buffers for each oscillator
	for i := 0; i != int(oscCount); i++ {
		oscOuts = append(oscOuts, in.EmptyBuf())
	}

	// Fill each oscOut with data from corresponding output
	return func(buf *audio.FloatBuffer) {
		for i := range inputs {
			inputs[i](oscOuts[i])
			for s := range buf.Data {
				buf.Data[s] += oscOuts[i].Data[s] / oscCount
			}
		}
	}
}

func (in *Instrument) WriteWAV(fileName string, out *audio.FloatBuffer) error {
	wavFile, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer wavFile.Close()

	enc := wav.NewEncoder(wavFile, out.PCMFormat().SampleRate, in.BitRate, out.PCMFormat().NumChannels, 1)
	if err := enc.Write(out.AsIntBuffer()); err != nil {
		return err
	}

	defer enc.Close()
	return nil
}
