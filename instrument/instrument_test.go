package instrument

import (
	. "github.com/IljaN/synth/generator"
	"reflect"
	"testing"
)

func TestOscCanBeModifiedByReferenceAfterAdd(t *testing.T) {
	p := New(44100, 16, 50000)
	ref := p.AddOscillator(WaveSaw, 440)
	ref.DcOffset = 9001

	if p.oscillators[0].DcOffset != 9001 || p.Oscillators()[0].DcOffset != 9001 {
		t.Errorf("Internal oscillator state was not changed when ref was modified")
	}
}

func TestAddRemoveOsc(t *testing.T) {
	p := New(44100, 16, 50000)

	_, osc2, _ := p.AddOscillator(WaveSqr, 440), p.AddOscillator(WaveSine, 330), p.AddOscillator(WaveSaw, 550)

	if exp, got := 3, len(p.Oscillators()); exp != got {
		t.Errorf("Expected %v oscillators got %v", exp, got)
	}

	osc2.Amplitude = 1337
	p.RemoveOscillator(osc2)

	if exp, got := 2, len(p.Oscillators()); exp != got {
		t.Errorf("Expected %v oscillators after remove got %v", exp, got)
	}

	for _, o := range p.Oscillators() {
		sf1 := reflect.ValueOf(o.Shape)
		sf2 := reflect.ValueOf(Sine)

		if o.Amplitude == 1337 || sf1.Pointer() == sf2.Pointer() {
			t.Errorf("Oscillator 2 was not removed")
		}
	}
}
