package bignum

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApproximation(t *testing.T) {
	sigmoid := func(x *big.Float) (y *big.Float) {
		z := new(big.Float).Set(x)
		z.Neg(z)
		z = Exp(z)
		z.Add(z, NewFloat(1, x.Prec()))
		y = NewFloat(1, x.Prec())
		y.Quo(y, z)
		return
	}

	prec := uint(128)

	t.Run("Chebyshev", func(t *testing.T) {

		interval := Interval{
			Nodes: 47,
			A:     *NewFloat(-4, prec),
			B:     *NewFloat(4, prec),
		}

		poly := ChebyshevApproximation(sigmoid, interval)

		xBig := NewFloat(1.4142135623730951, prec)

		y0, _ := sigmoid(xBig).Float64()
		y1, _ := poly.Evaluate(xBig)[0].Float64()

		require.InDelta(t, y0, y1, 1e-15)
	})

	t.Run("MultiIntervalMinimaxRemez", func(t *testing.T) {

		scanStep := NewFloat(1, prec)
		scanStep.Quo(scanStep, NewFloat(32, prec))

		k := 2.0

		intervals := []Interval{
			{A: *NewFloat(0, prec), B: *NewFloat(256/k, prec), Nodes: 36},
		}

		scaling := new(big.Float).Mul(Pi(prec), NewFloat(2/(8*k), prec))

		f := func(x *big.Float) (y *big.Float) {
			y = new(big.Float).Set(x)
			y.Mul(y, scaling)
			return Cos(y)
		}

		params := RemezParameters{
			Function:        f,
			Basis:           Chebyshev,
			Intervals:       intervals,
			ScanStep:        scanStep,
			Prec:            prec,
			OptimalScanStep: true,
		}

		r := NewRemez(params)
		r.Approximate(200, 1e-15)
		r.ShowCoeffs(50)
		r.ShowError(50)
	})
}
