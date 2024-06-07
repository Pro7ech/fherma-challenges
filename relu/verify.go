package main

import (
	"app/utils"
	"flag"
	"fmt"
	"log"
	"math"

	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
)

func main() {
	ccFile := flag.String("cc", "", "")
	skFile := flag.String("sk", "", "")
	inputFile := flag.String("input", "", "")
	outputFile := flag.String("output", "", "")

	flag.Parse()

	params := new(hefloat.Parameters)
	if err := utils.Deserialize(params, *ccFile); err != nil {
		log.Fatalf(err.Error())
	}

	sk := new(rlwe.SecretKey)
	if err := utils.Deserialize(sk, *skFile); err != nil {
		log.Fatalf(err.Error())
	}

	in := new(rlwe.Ciphertext)
	if err := utils.Deserialize(in, *inputFile); err != nil {
		log.Fatalf(err.Error())
	}

	out := new(rlwe.Ciphertext)
	if err := utils.Deserialize(out, *outputFile); err != nil {
		log.Fatalf(err.Error())
	}

	dec := rlwe.NewDecryptor(*params, sk)
	ecd := hefloat.NewEncoder(*params)

	have := make([]float64, out.Slots())
	if err := ecd.Decode(dec.DecryptNew(out), have); err != nil {
		log.Fatalf("%T.Decode: %s", ecd, err.Error())
	}

	want := make([]float64, in.Slots())
	if err := ecd.Decode(dec.DecryptNew(in), want); err != nil {
		log.Fatalf("%T.Decode: %s", ecd, err.Error())
	}

	for i := range have {
		if want[i] < 0 {
			want[i] = 0
		}
	}

	fmt.Println(have[:4])
	fmt.Println(want[:4])

	fmt.Println(hefloat.GetPrecisionStats(*params, ecd, nil, have, want, 0, false).String())

	var cnt int
	for i := range have {
		if math.Abs(have[i]-want[i]) <= 0.001 {
			cnt++
		}
	}

	fmt.Println(float64(cnt) / float64(len(want)) * 100)
	fmt.Println(have[8190:8194])
}
