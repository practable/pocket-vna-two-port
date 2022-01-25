package calibration

import (
	"errors"
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

func makeSParam(freq []uint64, real, imag []float64) ([]pocket.SParam, error) {

	if len(freq) != len(real) || len(freq) != len(imag) {
		return []pocket.SParam{}, errors.New("Freq/real/imag inconsistent length")
	}

	pa := []pocket.SParam{}

	for i := range freq {
		p := pocket.SParam{
			Freq: freq[i],
			S11: pocket.Complex{
				Imag: imag[i],
				Real: real[i],
			},
		}
		pa = append(pa, p)
	}

	return pa, nil
}

func TestMakeSParam(t *testing.T) {

	real := []float64{0.31, 0.32, -0.33}
	imag := []float64{0.41, -0.42, 0.45}
	freq := []uint64{1000, 2000, 3000}

	input, err := makeSParam(freq, real, imag)

	assert.NoError(t, err)

	expected_array := ComplexArray{
		Real: real,
		Imag: imag,
	}

	actual_freq, actual_array := PocketToCalibration(input)

	assert.Equal(t, freq, actual_freq)
	assert.Equal(t, expected_array, actual_array)

}

func TestMakeOnePort(t *testing.T) {

	freq := []uint64{1000, 2000, 3000}
	short_real := []float64{0.31, 0.32, -0.33}
	short_imag := []float64{0.41, -0.42, 0.45}
	open_real := []float64{1.31, 1.32, -1.33}
	open_imag := []float64{1.41, -1.42, 1.45}
	load_real := []float64{2.31, 2.32, -2.33}
	load_imag := []float64{2.41, -2.42, 2.45}
	dut_real := []float64{3.31, 3.32, -3.33}
	dut_imag := []float64{3.41, -3.42, 3.45}

	pshort, err := makeSParam(freq, short_real, short_imag)
	assert.NoError(t, err)
	popen, err := makeSParam(freq, open_real, open_imag)
	assert.NoError(t, err)
	pload, err := makeSParam(freq, load_real, load_imag)
	assert.NoError(t, err)
	pdut, err := makeSParam(freq, dut_real, dut_imag)
	assert.NoError(t, err)

	op, err := MakeOnePort(pshort, popen, pload, pdut)

	assert.NoError(t, err)

	assert.Equal(t, freq, op.Freq)

	assert.Equal(t, short_real, op.Short.Real)
	assert.Equal(t, short_imag, op.Short.Imag)
	assert.Equal(t, open_real, op.Open.Real)
	assert.Equal(t, open_imag, op.Open.Imag)
	assert.Equal(t, load_real, op.Load.Real)
	assert.Equal(t, load_imag, op.Load.Imag)
	assert.Equal(t, dut_real, op.DUT.Real)
	assert.Equal(t, dut_imag, op.DUT.Imag)

	// freq length does not match data length
	freq = []uint64{1000, 2000}
	dut_real = []float64{3.31, 3.32}
	dut_imag = []float64{3.41, -3.42}
	pdut, err = makeSParam(freq, dut_real, dut_imag)
	assert.NoError(t, err)
	op, err = MakeOnePort(pshort, popen, pload, pdut)
	assert.Error(t, err)

	// data lengths are different
	freq = []uint64{1000, 2000, 3000, 4000}
	dut_real = []float64{3.31, 3.32, -3.33, 5.66}
	dut_imag = []float64{3.41, -3.42, 3.45, 6.77}
	pdut, err = makeSParam(freq, dut_real, dut_imag)
	assert.NoError(t, err)
	op, err = MakeOnePort(pshort, popen, pload, pdut)
	assert.Error(t, err)

}
