package main

import (
	"github.com/IljaN/synth/encoding"
	"github.com/IljaN/synth/filter"
	"github.com/IljaN/synth/generator"
	inst "github.com/IljaN/synth/instrument"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	_ = os.Remove("bd.wav")

	s := inst.New(48000, 16, 500000)
	out := s.Master(
		inst.Chain(
			s.Mix(inst.OscOut(s.NewOsc(generator.Sine, 230)), inst.OscOut(s.NewOsc(generator.WhiteNoise, 330))),
			filter.NewLPF(&filter.LPF{Cutoff: 0.4353535}),
			filter.NewDelayFilter(filter.DelayFilterParams(0.525, 0.2325, 0.03532)),
		),
	)

	if err := encoding.WriteWAV("bd.wav", out, 16); err != nil {
		log.Fatal(err)
	}

	ffplay("bd.wav", 48000)
}

// plays and visualizes the generated sound with ffplay
func ffplay(fn string, sampleRate int) {
	ffplayExecutable, _ := exec.LookPath("ffplay")
	ffplayCmd := &exec.Cmd{
		Path:   ffplayExecutable,
		Args:   []string{ffplayExecutable, "-showmode", "1", "-ar", strconv.Itoa(sampleRate), fn},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := ffplayCmd.Start(); err != nil {
		panic(err)
	}

	ffplayCmd.Wait()
}
