package main

import (
	"github.com/IljaN/synth/filter"
	"github.com/IljaN/synth/generator"
	inst "github.com/IljaN/synth/instrument"
	"os"
	"os/exec"
)

func main() {
	_ = os.Remove("bd.wav")
	s := inst.New(48000, 16, 500000)
	s.WriteWAV("bd.wav", s.Out(inst.Channel(s.Mix(
		inst.OscOutput(s.NewOsc(generator.Sine, 230, 48000)),
		inst.OscOutput(s.NewOsc(generator.WhiteNoise, 330, 48000)),
	),
		filter.NewDelayFilter(0.6, 0.3, 0.5).Filter,
		filter.NewLPF(0.00110).Filter,
		filter.NewFlangerFilter(0.54, 0.6, 0.9).Filter,
	)))

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
