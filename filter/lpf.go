package filter

import (
	"github.com/go-audio/audio"
	"math"
)

type LPF struct {
	previousLeft  float64
	previousRight float64
	Cutoff        float64
}

func NewLPF(cutoff float64) *LPF {
	return &LPF{
		previousLeft:  0.0,
		previousRight: 0.0,
		Cutoff:        cutoff,
	}
}

func (f *LPF) Filter(buf *audio.FloatBuffer) {
	n := len(buf.Data)
	isStereo := buf.Format.NumChannels == 2
	if isStereo {
		n = n / 2
	}
	waveLength := float64(n) / float64(buf.Format.SampleRate)
	rc := 1.0 / (2 * math.Pi * f.Cutoff)
	alpha := waveLength / (rc + waveLength)
	for i := 0; i < n; i++ {
		if isStereo {
			buf.Data[i*2] = f.previousLeft + alpha*(buf.Data[i*2]-f.previousLeft)
			buf.Data[i*2+1] = f.previousRight + alpha*(buf.Data[i*2+1]-f.previousRight)
			f.previousLeft = buf.Data[i*2]
			f.previousRight = buf.Data[i*2+1]
		} else {
			buf.Data[i] = f.previousLeft + alpha*(buf.Data[i]-f.previousLeft)
			f.previousLeft = buf.Data[i]
		}
	}

}
