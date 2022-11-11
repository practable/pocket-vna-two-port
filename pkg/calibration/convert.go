package calibration

import (
	"errors"

	"github.com/timdrysdale/go-pocketvna/pkg/pocket"
)

func MakeOnePort(pshort, popen, pload, pdut []pocket.SParam) (Command, error) {

	freq, cshort := PocketToCalibration(pshort)

	if len(freq) != len(cshort.Real) {
		return Command{}, errors.New("data set mismatch")
	}

	tmp, copen := PocketToCalibration(popen)

	if len(freq) != len(tmp) || len(freq) != len(copen.Real) {
		return Command{}, errors.New("data set mismatch")
	}

	tmp, cload := PocketToCalibration(pload)

	if len(freq) != len(tmp) || len(freq) != len(cload.Real) {
		return Command{}, errors.New("data set mismatch")
	}

	tmp, cdut := PocketToCalibration(pdut)

	if len(freq) != len(tmp) || len(freq) != len(cdut.Real) {
		return Command{}, errors.New("data set mismatch")
	}

	return Command{
		Command: "oneport",
		Freq:    freq,
		Short:   cshort,
		Open:    copen,
		Load:    cload,
		DUT:     cdut,
	}, nil

}

// PocketToCalibration changes the array structure to better suit our
// usages of the results (e.g. calibration calculations, and presentation
// in the user interface both need a frequency list so let's prepare that now
// for efficiency)
func PocketToCalibration(p []pocket.SParam) ([]uint64, SParam) {

	// we'll use append rather than assuming max-length array
	var freq []uint64
	var s11_real, s11_imag, s12_real, s12_imag, s21_real, s21_imag, s22_real, s22_imag []float64

	for _, param := range p {
		freq = append(freq, param.Freq)
		s11_real = append(real, param.S11.Real)
		s11_imag = append(imag, param.S11.Imag)
		s12_real = append(real, param.S12.Real)
		s12_imag = append(imag, param.S12.Imag)
		s21_real = append(real, param.S21.Real)
		s21_imag = append(imag, param.S21.Imag)
		s22_real = append(real, param.S22.Real)
		s22_imag = append(imag, param.S22.Imag)
	}

	sp := SParam{
		S11: ComplexArray{
			Real: s11_real,
			Imag: s11_imag,
		},
		S12: ComplexArray{
			Real: s12_real,
			Imag: s12_imag,
		},
		S21: ComplexArray{
			Real: s21_real,
			Imag: s21_imag,
		},
		S22: ComplexArray{
			Real: s22_real,
			Imag: s22_imag,
		},
	}
	return freq, sp

}

func PocketToResult(p []pocket.SParam) Result {
	freq, ca := PocketToCalibration(p)
	return Result{
		Freq: freq,
		S11:  sp.S11,
		S12:  sp.S12,
		S21:  sp.S21,
		S22:  sp.S22,
	}
}

func CalibrationToPocket(result Result) ([]pocket.SParam, error) {

	pa := []pocket.SParam{}

	if len(result.Freq) != len(result.S11.Real) {
		return pa, errors.New("Freq and S11 are different lengths")
	}

	for i, freq := range result.Freq {
		p := pocket.SParam{
			Freq: freq,
			S11: pocket.Complex{
				Real: result.S11.Real[i],
				Imag: result.S11.Imag[i],
			},
		}
		pa = append(pa, p)
	}

	return pa, nil
}
