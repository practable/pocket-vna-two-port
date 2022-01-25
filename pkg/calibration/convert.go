package calibration

import "github.com/timdrysdale/go-pocketvna/pkg/pocket"

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
