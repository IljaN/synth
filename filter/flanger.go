package filter

import (
	"github.com/go-audio/audio"
	"math"
)

type FlangerFilter struct {
	Time    float64
	Factor  float64
	LFORate float64

	LeftDelayed  []float64
	RightDelayed []float64
	Phase        int
}

func NewFlangerFilter(time, factor, rate float64) *FlangerFilter {
	return &FlangerFilter{
		Time:         time,
		Factor:       factor,
		LFORate:      rate,
		LeftDelayed:  nil,
		RightDelayed: nil,
	}
}

func (f *FlangerFilter) Filter(buf *audio.FloatBuffer) error {
	sampleRate := buf.Format.SampleRate
	isStereo := buf.Format.NumChannels == 2
	time := int(float64(sampleRate) * f.Time)
	if f.LeftDelayed == nil {
		f.LeftDelayed = make([]float64, time)
		f.RightDelayed = make([]float64, time)
	}

	n := len(buf.Data)
	if isStereo {
		n = n / 2
	}
	stepSize := (f.LFORate * math.Pi) / float64(sampleRate)
	for i := 0; i < n; i++ {

		ix := i
		if isStereo {
			ix *= 2
		}

		currentDelay := time - int(math.Ceil(float64(time)*math.Abs(math.Sin(float64(f.Phase)*stepSize))))
		delayedIx := f.Phase - currentDelay
		if delayedIx < 0 {
			delayedIx += time
		}

		f.LeftDelayed[f.Phase] = buf.Data[ix]
		if isStereo {
			f.RightDelayed[f.Phase] = buf.Data[ix+1]
		}

		buf.Data[ix] = (1.0-f.Factor)*buf.Data[ix] + f.Factor*f.LeftDelayed[delayedIx]
		if isStereo {
			buf.Data[ix+1] = (1.0-f.Factor)*buf.Data[ix+1] + f.Factor*f.RightDelayed[delayedIx]
		}

		f.Phase += 1
		if f.Phase >= len(f.LeftDelayed) {
			f.Phase = 0
		}
	}

	return nil
}
