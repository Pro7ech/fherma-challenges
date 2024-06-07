package solution

import (
	"math"
	"math/big"

	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
	"github.com/tuneinsight/lattigo/v5/utils/bignum"
)

var coeffs = [][]string{
	{"0", "0.750602769761", "0", "-0.255359946795", "0", "0.160747549652", "0", "-0.479562684008"},
	{"0", "1.230920518721", "0", "-0.397254039291", "0", "0.223294137470", "0", "-0.144110311154", "0", "0.097569608332", "0", "-0.066606086948", "0", "0.044658190057", "0", "-0.046554908569"},
	{"0", "1.236492071810", "0", "-0.325598760814", "0", "0.120926543606", "0", "-0.041146355130", "0", "0.011358439319", "0", "-0.002320523292", "0", "0.000308521301", "0", "-0.000019936801"},
}

func SolveTestcase0(params *hefloat.Parameters, evk *rlwe.MemEvaluationKeySet, in *rlwe.Ciphertext) (out *rlwe.Ciphertext, err error) {

	/*
		hefloat.GenMinimaxCompositePolynomialForSign(
			256,
			6,
			12,
			[]int{
				7,
				15, // 5
				15, // 6
			})
	*/

	eval := hefloat.NewEvaluator(*params, evk)

	polys := hefloat.NewMinimaxCompositePolynomial(coeffs)

	CmpEval := hefloat.NewComparisonEvaluator(*params, eval, nil, polys)

	var step *rlwe.Ciphertext
	if step, err = CmpEval.Step(in); err != nil {
		return
	}

	if err = eval.MulRelin(in, step, in); err != nil {
		return
	}

	if err = eval.Rescale(in, in); err != nil {
		return
	}

	return in, nil
}

func SolveTestcase1(params *hefloat.Parameters, evk *rlwe.MemEvaluationKeySet, in *rlwe.Ciphertext) (out *rlwe.Ciphertext, err error) {

	var prec uint = 128

	scanStep := bignum.NewFloat(1, prec)
	scanStep.Quo(scanStep, bignum.NewFloat(32, prec))

	f := func(x *big.Float) (y *big.Float) {
		y = new(big.Float).Set(x)
		if y.Cmp(new(big.Float)) < 1 {
			y.SetFloat64(0)
		}
		return
	}

	a := 1.0
	b := 0.284723851

	intervals := []bignum.Interval{
		{A: *bignum.NewFloat(-a, prec), B: *bignum.NewFloat(-b, prec), Nodes: 5},
		{A: *bignum.NewFloat(b, prec), B: *bignum.NewFloat(a, prec), Nodes: 5},
	}

	r := bignum.NewRemez(bignum.RemezParameters{
		Function:        f,
		Basis:           bignum.Chebyshev,
		Intervals:       intervals,
		ScanStep:        scanStep,
		Prec:            prec,
		OptimalScanStep: true,
	})
	r.Approximate(200, 1e-15)
	r.ShowCoeffs(15)
	r.ShowError(15)

	coeffs := make([]float64, 9)
	for i := range coeffs {
		coeffs[i], _ = r.Coeffs[i].Float64()
	}

	eval := hefloat.NewEvaluator(*params, evk)

	degree := len(coeffs) - 1

	d := hefloat.NewPowerBasis(in, bignum.Chebyshev)
	d.GenPower(degree, false, eval)

	scaling := coeffs[degree]

	for i := range coeffs {
		coeffs[i] /= scaling
	}

	evalPoly := hefloat.NewPolynomialEvaluator(*params, eval)

	poly := hefloat.NewPolynomial(bignum.NewPolynomial(bignum.Chebyshev, coeffs[:8], [2]float64{-1, 1}))

	if out, err = evalPoly.Evaluate(in, poly, d.Value[8].Scale); err != nil {
		return
	}

	if err = eval.Add(out, d.Value[8], out); err != nil{
		return
	}

	out.Scale = out.Scale.Div(rlwe.NewScale(math.Abs(scaling) / a))
	
	if err = eval.Mul(out, -1, out); err != nil{
		return
	}

	return out, nil
}
