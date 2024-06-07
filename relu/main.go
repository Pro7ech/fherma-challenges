package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"

	"app/internal/solution"
	"app/utils"
)

func main() {
	var err error

	ccFile := flag.String("cc", "", "")
	evkFile := flag.String("key_eval", "", "")
	inputFile := flag.String("input", "", "")
	outputFile := flag.String("output", "", "")

	flag.Parse()

	params := new(hefloat.Parameters)
	if err = utils.Deserialize(params, *ccFile); err != nil {
		log.Fatalf(err.Error())
	}

	rlk := rlwe.RelinearizationKey{}
	if err = utils.Deserialize(&rlk, *evkFile); err != nil {
		log.Fatalf(err.Error())
	}

	in := new(rlwe.Ciphertext)
	if err = utils.Deserialize(in, *inputFile); err != nil {
		log.Fatalf(err.Error())
	}

	var out *rlwe.Ciphertext
	fmt.Println(*ccFile)

	if *ccFile == "./artifact/context_0" {
		if out, err = solution.SolveTestcase0(params, rlwe.NewMemEvaluationKeySet(&rlk), in); err != nil {
			log.Fatalf("solution.SolveTestcase0: %s", err.Error())
		}
	} else {
		if out, err = solution.SolveTestcase1(params, rlwe.NewMemEvaluationKeySet(&rlk), in); err != nil {
			log.Fatalf("solution.SolveTestcase0: %s", err.Error())
		}
	}

	utils.Serialize(out, *outputFile)
}
