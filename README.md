# synth

Dirty functional synthesizer DSL.

```golang
s := inst.New(48000, 16, 500000)
out := s.Master(
    // Mix two osc`s and pass the result trough a filter-chain
    inst.Chain(
        s.Mix(inst.OscOut(s.NewOsc(generator.Sine, 230)), inst.OscOut(s.NewOsc(generator.WhiteNoise, 330))),
        filter.NewLPF(0.00110).Filter,
        filter.NewDelayFilter(0.06, 0.03343, 0.05).Filter,
    ),
)

if err := encoding.WriteWAV("sound.wav", out, 16); err != nil {
    log.Fatal(err)
}

```

### Custom generator functions
```golang
s := inst.New(48000, 16, 500000)
out := s.Master(
    inst.OscOut(generator.NewOsc(func(x float64) float64 {
        return 1.0 / math.SqrtPhi * x
    }, 140, 16)),
)

if err := encoding.WriteWAV("sound.wav", out, 16); err != nil {
    log.Fatal(err)
}

// Visualize with ffplay
ffplay("sound.wav", 48000)
```

### Lazy higher-order transformation
```golang
inst.Chain(
    filter.NewDelayFilter(0.0533, 0.002, 0.002).Filter,
    filter.NewFlangerFilter(0.053,0.03,0.2).Filter,
    func(buf *audio.FloatBuffer) { /* Custom transform. */ })
```