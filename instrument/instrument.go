package instrument

import (
	"github.com/go-audio/audio"
	. "github.com/go-audio/generator"
	"github.com/go-audio/wav"
	"os"
)

type Instrument struct {
	oscillators  []*Osc
	filters      []Filter
	out          *audio.FloatBuffer
	SamplingRate int
	BitRate      int
	Duration     int
}

type Filter interface {
	Filter(buf *audio.FloatBuffer) error
}

func New(samplingRate int, bitRate int, duration int) *Instrument {
	p := &Instrument{
		oscillators:  make([]*Osc, 0),
		filters:      make([]Filter, 0),
		SamplingRate: samplingRate,
		BitRate:      bitRate,
		Duration:     duration,
	}

	p.out = p.emptyBuf()
	return p
}

func (in *Instrument) Render() error {
	oscCount := float64(len(in.oscillators))
	oscOuts := make([]*audio.FloatBuffer, 0, int(oscCount))

	// Pre allocate out buffers for each oscillator
	for i := 0; i != int(oscCount); i++ {
		oscOuts = append(oscOuts, in.emptyBuf())
	}

	// Fill each oscOut with data from corresponding osc
	for i := range in.oscillators {
		_ = in.oscillators[i].Fill(oscOuts[i])
		for s := range in.out.Data {
			in.out.Data[s] += oscOuts[i].Data[s] / oscCount
		}
	}

	// Apply filters on sum
	for i := range in.filters {
		if err := in.filters[i].Filter(in.out); err != nil {
			return err
		}
	}

	return nil
}

func (in *Instrument) WriteWAV(fileName string) error {
	if err := in.Render(); err != nil {
		return err
	}

	wavFile, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer wavFile.Close()

	enc := wav.NewEncoder(wavFile, in.out.PCMFormat().SampleRate, in.BitRate, in.out.PCMFormat().NumChannels, 1)
	if err := enc.Write(in.out.AsIntBuffer()); err != nil {
		return err
	}

	defer enc.Close()
	return nil
}

func (in *Instrument) Write(buf *audio.FloatBuffer, enc wav.Encoder) error {
	if err := enc.Write(buf.AsIntBuffer()); err != nil {
		return err
	}
	return nil
}

func (in *Instrument) emptyBuf() *audio.FloatBuffer {
	return &audio.FloatBuffer{
		Data:   make([]float64, in.Duration*2),
		Format: audio.FormatStereo48000,
	}
}

func (in *Instrument) Oscillators() []*Osc {
	return in.oscillators
}

func (in *Instrument) AddOscillator(shape WaveType, freqHz float64) *Osc {
	osc := NewOsc(shape, freqHz, in.SamplingRate)
	osc.Amplitude = float64(audio.IntMaxSignedValue(in.BitRate))
	in.oscillators = append(in.oscillators, osc)

	return osc
}

func (in *Instrument) RemoveOscillator(osc *Osc) (ok bool) {
	for i := range in.oscillators {
		if in.oscillators[i] == osc {
			in.oscillators = append(in.oscillators[:i], in.oscillators[i+1:]...)
			return true
		}
	}

	return false
}

func (in *Instrument) SetOscillators(o ...*Osc) []*Osc {
	in.oscillators = o
	return in.oscillators
}

func (in *Instrument) Filters() []Filter {
	return in.filters
}

func (in *Instrument) AddFilter(f Filter) {
	in.filters = append(in.filters, f)
}

func (in *Instrument) RemoveFilter(f Filter) (ok bool) {
	for i := range in.filters {
		if in.filters[i] == f {
			in.filters = append(in.filters[:i], in.filters[i+1:]...)
			return true
		}
	}

	return false
}

func (in *Instrument) SetFilters(chain ...Filter) {
	in.filters = chain
}
