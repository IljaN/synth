package main

import (
	"github.com/IljaN/synth/filter"
	"github.com/IljaN/synth/generator"
	inst "github.com/IljaN/synth/instrument"
	"github.com/go-audio/audio"
	"os"
	"os/exec"
)

func main() {
	_ = os.Remove("bd.wav")
	p := inst.New(48000, 16, 500000)
	out := p.Out(inst.Channel(
		p.Mix(func(buf *audio.FloatBuffer) {
			osc1 := generator.NewOsc(generator.Square, 230, 48000)
			osc1.Amplitude = float64(audio.IntMaxSignedValue(16))
			osc1.Fill(buf)
		}, func(buf *audio.FloatBuffer) {
			osc2 := generator.NewOsc(generator.WaveWhiteNoise, 330, 48000)
			osc2.Amplitude = float64(audio.IntMaxSignedValue(16))
			osc2.Fill(buf)
		}),
		filter.NewDelayFilter(0.6, 0.3, 0.5).Filter,
		filter.NewLPF(0.00110).Filter,
		filter.NewFlangerFilter(0.54, 0.6, 0.9).Filter,
	))

	p.WriteWAV("bd.wav", out)

	ffplayExecutable, _ := exec.LookPath("ffplay")
	ffplayCmd := &exec.Cmd{
		Path:   ffplayExecutable,
		Args:   []string{ffplayExecutable, "-showmode", "1", "-ar", "44100", "bd.wav"},
		Env:    nil,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := ffplayCmd.Start(); err != nil {
		panic(err)
	}

	ffplayCmd.Wait()

}
