package generator

import (
	"math"

	"github.com/go-audio/audio"
)

// Osc is an oscillator
type Osc struct {
	Shape     WaveFunc
	Amplitude float64
	DcOffset  float64
	Freq      float64
	// SampleRate
	Fs                int
	PhaseOffset       float64
	CurrentPhaseAngle float64
	phaseAngleIncr    float64
	// currentSample allows us to track where we are at in the signal life
	// and setup an envelope accordingly
	currentSample int
	// ADSR
	attackInSamples int
}

// NewOsc returns a new oscillator, note that if you change the phase offset of the returned osc,
// you also need to set the CurrentPhaseAngle
func NewOsc(shape WaveFunc, hz float64, fs int) *Osc {
	return &Osc{Shape: shape, Amplitude: 1, Freq: hz, Fs: fs, phaseAngleIncr: ((hz * TwoPi) / float64(fs))}
}

// Reset sets the oscillator back to its starting state
func (o *Osc) Reset() {
	o.phaseAngleIncr = ((o.Freq * TwoPi) / float64(o.Fs))
	o.currentSample = 0
}

// SetFreq updates the oscillator frequency
func (o *Osc) SetFreq(hz float64) {
	if o.Freq != hz {
		o.Freq = hz
		o.phaseAngleIncr = ((hz * TwoPi) / float64(o.Fs))
	}
}

// SetAttackInMs sets the duration for the oscillator to be at full amplitude
// after it starts.
func (o *Osc) SetAttackInMs(ms int) {
	if o == nil {
		return
	}
	if ms <= 0 {
		o.attackInSamples = 0
		return
	}
	o.attackInSamples = int(float32(o.Fs) / (1000.0 / float32(ms)))
}

// Signal uses the osc to generate a discreet signal
func (o *Osc) Signal(length int) []float64 {
	output := make([]float64, length)
	for i := 0; i < length; i++ {
		output[i] = o.Sample()
	}
	return output
}

// Fill fills up the pass audio Buffer with the output of the oscillator.
func (o *Osc) Fill(buf *audio.FloatBuffer) error {
	if o == nil {
		return nil
	}
	numChans := 1
	if f := buf.Format; f != nil {
		numChans = f.NumChannels
	}
	frameCount := buf.NumFrames()
	var sample float64
	for i := 0; i < frameCount; i++ {
		sample = o.Sample()
		for j := 0; j < numChans; j++ {
			buf.Data[i*numChans+j] = sample
		}
	}
	return nil
}

// Sample returns the next sample generated by the oscillator
func (o *Osc) Sample() (output float64) {
	if o == nil {
		return
	}
	o.currentSample++
	if o.CurrentPhaseAngle < -math.Pi {
		o.CurrentPhaseAngle += TwoPi
	} else if o.CurrentPhaseAngle > math.Pi {
		o.CurrentPhaseAngle -= TwoPi
	}

	var amp float64
	if o.attackInSamples > o.currentSample {
		// linear fade in
		amp = float64(o.currentSample) * (o.Amplitude / float64(o.attackInSamples))
	} else {
		amp = o.Amplitude
	}

	output = amp*o.Shape(o.CurrentPhaseAngle) + o.DcOffset
	o.CurrentPhaseAngle += o.phaseAngleIncr
	return output
}
