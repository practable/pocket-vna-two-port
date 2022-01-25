package calibration

import (
	"errors"

	"github.com/timdrysdale/go-pocketvna/pkg/pocket"
)

func MakeOnePort(pshort, popen, pload, pdut []pocket.SParam) (Command, error) {

	freq, cshort := PocketToCalibration(pshort)

	tmp, copen := PocketToCalibration(popen)

	if len(freq) != len(tmp) {
		return Command{}, errors.New("data set mismatch")
	}

	tmp, cload := PocketToCalibration(pload)

	if len(freq) != len(tmp) {
		return Command{}, errors.New("data set mismatch")
	}

	tmp, cdut := PocketToCalibration(pdut)

	if len(freq) != len(tmp) {
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

func PocketToCalibration(p []pocket.SParam) ([]uint64, ComplexArray) {

	// we'll use append rather than assuming max-length array
	var freq []uint64
	var real, imag []float64

	for _, param := range p {
		freq = append(freq, param.Freq)
		real = append(real, param.S11.Real)
		imag = append(imag, param.S11.Imag)
	}

	ca := ComplexArray{
		Real: real,
		Imag: imag,
	}

	return freq, ca

}

func CalibrationToPocket() {
}
