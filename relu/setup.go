package main

import (
	"app/utils"
	"flag"
	"log"

	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
)

func main() {
	ccFile := flag.String("cc", "", "")
	skFile := flag.String("sk", "", "")
	evalFile := flag.String("key_eval", "", "")
	inputFile := flag.String("input", "", "")

	flag.Parse()

	params := hefloat.Parameters{}
	if err := utils.Deserialize(&params, *ccFile); err != nil {
		log.Fatalf(err.Error())
	}

	kgen := rlwe.NewKeyGenerator(params)

	sk := kgen.GenSecretKeyNew()

	ecd := hefloat.NewEncoder(params)

	enc := rlwe.NewEncryptor(params, sk)

	rlk := kgen.GenRelinearizationKeyNew(sk)

	values := make([]float64, params.MaxSlots())
	for i := range values {
		values[i] = 2*(float64(i)/float64(params.MaxSlots())) - 1
	}

	pt := hefloat.NewPlaintext(params, params.MaxLevel())

	if err := ecd.Encode(values, pt); err != nil {
		log.Fatalf(err.Error())
	}

	input, err := enc.EncryptNew(pt)

	if err != nil {
		log.Fatalf(err.Error())
	}

	if err := utils.Serialize(params, *ccFile); err != nil {
		log.Fatalf(err.Error())
	}

	if err := utils.Serialize(sk, *skFile); err != nil {
		log.Fatalf(err.Error())
	}

	if err := utils.Serialize(rlk, *evalFile); err != nil {
		log.Fatalf(err.Error())
	}

	if err := utils.Serialize(input, *inputFile); err != nil {
		log.Fatalf(err.Error())
	}
}
