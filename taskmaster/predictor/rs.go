package predictor

import (
	"math/rand"
)

type RS struct {
	// slice of functions
	functions []string
	// map of functions to their fnParameters
	fnSet map[string]map[string]string
}

func NewRS() *RS {
	var rs = new(RS)
	rs.functions = make([]string, 0, 16)
	rs.fnSet = make(map[string]map[string]string)
	return rs
}

func (rs *RS) Update(info map[string]any) {
	// map should contain
	//  1. fnRequest: FnRequest struct
	var fnRequest FnRequest = info["fnRequest"].(FnRequest)

	// if function is not in fnSet, add it and update the function slice
	_, contains := rs.fnSet[fnRequest.FnName]
	if !contains {
		rs.fnSet[fnRequest.FnName] = fnRequest.FnParameters
		rs.functions = append(rs.functions, fnRequest.FnName)
	}
}

// Predict returns a random function from the function slice
func (rs *RS) Predict() *Prediction {
	var n int = len(rs.functions)
	// This means there is nothing in the slice
	if n == 0 {
		return NilPrediction
	}

	var i int = rand.Intn(n)
	var name string = rs.functions[i]

	return &Prediction{
		FnName:       name,
		FnParameters: rs.fnSet[name],
	}
}
