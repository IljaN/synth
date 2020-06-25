package main

import (
	"github.com/IljaN/synth/instrument"
	"github.com/go-audio/generator"
	"os"
	"os/exec"
)

func main() {

	_ = os.Remove("bd.wav")
	p := instrument.New(48000, 16, 50000)
	p.AddOscillator(generator.WaveSine, 330)
	p.AddOscillator(generator.WaveSaw, 220)

	/*
		p.SetFilters(
			filter.NewDelayFilter(0.3150,0.6, 0.3200),
			filter.NewLPF(0.00015),
		)

	*/

	p.WriteWAV("bd.wav")
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
