package predictor

// "container/list"
// "gopkg.in/yaml.v2"

// "os"

type MFE struct {
	// map of functions with key as function and count of invocation times
	functionCounts map[string]int64
	// keep track of max function
	currentName     string
	currentMaxCount int64
	// map of functions to their fnParameters
	fnSet map[string]map[string]string
}

func NewMFE() *MFE {
	var mfe = new(MFE)
	mfe.functionCounts = make(map[string]int64)
	mfe.fnSet = make(map[string]map[string]string)
	mfe.currentMaxCount = 0
	return mfe
}

func (mfe *MFE) Update(info map[string]any) {
	// map should contain
	//  1. fnRequest: FnRequest struct
	var fnRequest FnRequest = info["fnRequest"].(FnRequest)
	// log.Printf("%+v\n", fnRequest)
	// MFE algorithm
	// First check if it exists in the list
	_, contains := mfe.functionCounts[fnRequest.FnName]
	// log.Printf("Contains = %v\n", contains)
	if !contains {
		mfe.functionCounts[fnRequest.FnName] = 1
	} else {
		mfe.functionCounts[fnRequest.FnName] += 1
	}
	// update currentName and currentMaxCount
	if mfe.currentMaxCount < mfe.functionCounts[fnRequest.FnName] {
		mfe.currentMaxCount = mfe.functionCounts[fnRequest.FnName]
		mfe.currentName = fnRequest.FnName
	}
}

// Predict returns the last item in the internal linked list
func (mfe *MFE) Predict() *Prediction {
	// This means maxcount has never been updated
	if mfe.currentMaxCount == 0 {
		return NilPrediction
	}

	return &Prediction{
		FnName:       mfe.currentName,
		FnParameters: mfe.fnSet[mfe.currentName],
	}
}
