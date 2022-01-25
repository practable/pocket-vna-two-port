package calibration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/go-pocketvna/pkg/pocket"
)

func TestConvertPocketToCalibration(t *testing.T) {

	real := []float64{0.31, 0.32, -0.33}
	imag := []float64{0.41, -0.42, 0.45}
	freq := []uint64{1000, 2000, 3000}

	input := []pocket.SParam{
		pocket.SParam{
			Freq: freq[0],
			S11: pocket.Complex{
				Real: real[0],
				Imag: imag[0],
			},
		},
		pocket.SParam{
			Freq: freq[1],
			S11: pocket.Complex{
				Real: real[1],
				Imag: imag[1],
			},
		},
		pocket.SParam{
			Freq: freq[2],
			S11: pocket.Complex{
				Real: real[2],
				Imag: imag[2],
			},
		},
	}

	expected_array := ComplexArray{
		Real: real,
		Imag: imag,
	}

	actual_freq, actual_array := PocketToCalibration(input)

	assert.Equal(t, freq, actual_freq)
	assert.Equal(t, expected_array, actual_array)

}
