package generator

import "testing"

func TestOsc_Signal(t *testing.T) {
	osc := NewOsc(Sine, 440, 44100)
	if osc.CurrentPhaseAngle != 0 {
		t.Fatalf("expected the current phase to be zero")
	}
	if osc.phaseAngleIncr != 0.06268937721449021 {
		t.Fatalf("Wrong phase angle increment")
	}
	sample := osc.Sample()
	if phase := osc.CurrentPhaseAngle; phase != 0.06268937721449021 {
		t.Fatalf("wrong phase angle: %f, expected 0.06268937721449021", phase)
	}
	if sample != 0.0 {
		t.Fatalf("wrong first sample: %f expected 0.0", sample)
	}
	signal := osc.Signal(19)
	expected := []float64{0.062001866171879985, 0.12406665714043713, 0.1858716671561314, 0.24710788950751222, 0.3074800165212184,
		0.3667064395619785, 0.42451924903261046, 0.4806642343740216, 0.5349008840652089, 0.5870023856232587, 0.6367556256033469,
		0.6839611895987389, 0.7284333622407899, 0.7700001271989437, 0.8085031671807348, 0.8437978639317864, 0.8757532982358108,
		0.9042522499146113, 0.9291911978280791}

	for i, s := range signal {
		if !nearlyEqual(s, expected[i], 0.000001) {
			t.Logf("sample %d didn't match, expected: %f got %f\n", i, expected[i], s)
			t.Fail()
		}
	}

	osc = NewOsc(Sine, 400, 1000)
	signal = osc.Signal(100)

	expected = []float64{0, 0.5881600000000001, -0.9513600000000001, 0.9513600000000001, -0.5881600000000001, 0, 0.5881600000000001, -0.9513600000000001, 0.9513600000000001, -0.5881600000000001, 0, 0.5881600000000001, -0.9513600000000001, 0.9513600000000001, -0.5881600000000001, 0, 0.5881600000000001, -0.9513600000000001, 0.9513600000000001, -0.5881600000000001, 0, 0.5881600000000001, -0.9513600000000001, 0.9513600000000001, -0.5881600000000001, 0, 0.5881600000000001, -0.9513600000000001, 0.9513600000000001, -0.5881600000000001, 0, 0.5881600000000001, -0.9513600000000001, 0.9513600000000001, -0.5881600000000001, 0, 0.5881600000000001, -0.9513600000000001, 0.9513600000000001, -0.5881600000000001, 0, 0.5881600000000001, -0.9513600000000001, 0.9513600000000001, -0.5881600000000001, 0, 0.5881600000000001, -0.9513600000000001, 0.9513600000000001, -0.5881600000000001, 0, 0.5881600000000001, -0.9513600000000001, 0.9513600000000001, -0.5881600000000001, 0, 0.5881600000000001, -0.9513600000000001, 0.9513600000000001, -0.5881600000000001, 0, 0.5881600000000001, -0.9513600000000001, 0.9513600000000001, -0.5881600000000001, 0, 0.5881600000000001, -0.9513600000000001, 0.9513600000000001, -0.5881600000000001, 0, 0.5881600000000001, -0.9513600000000001, 0.9513600000000001, -0.5881600000000001, 0, 0.5881600000000001, -0.9513600000000001, 0.9513600000000001, -0.5881600000000001, 0, 0.5881600000000001, -0.9513600000000001, 0.9513600000000001, -0.5881600000000001, 0, 0.5881600000000001, -0.9513600000000001, 0.9513600000000001, -0.5881600000000001, 0, 0.5881600000000001, -0.9513600000000001, 0.9513600000000001, -0.5881600000000001, 0, 0.5881600000000001, -0.9513600000000001, 0.9513600000000001, -0.5881600000000001}
	for i, s := range signal {
		if !nearlyEqual(s, expected[i], 0.00000001) {
			t.Logf("sample %d didn't match, expected: %f got %f\n", i, expected[i], s)
			t.Fail()
		}
	}
}
