package predictor

// FnRequest represents a function request
type FnRequest struct {
	FnName       string
	FnParameters map[string]string
}

// EqualsFnRequest is a helper method for determining if two FnRequests are the same
func EqualsFnRequest(a, b FnRequest) bool {
	// It delegates to a simple equality check on the FnName since FnParameters isn't important
	return a.FnName == b.FnName
}

// Prediction is the output made by the Predictor
type Prediction struct {
	// FnName is the name of the function to be called
	FnName string
	// FnParameters are the parameters to call the function with
	FnParameters map[string]string
}

// NilPrediction is returned if there is no function to predict yet
var NilPrediction = &Prediction{FnName: "", FnParameters: nil}

// Predictor is the interface that taskmaster uses to decide which function to ping to keep alive.
type Predictor interface {
	// Update updates the Predictor's internal state
	// everytime a serverless function is invoked for real
	Update(info map[string]any)
	// Predict returns the next serverless function to ping
	// if the serverless function requires parameters, they will be randomly generated
	Predict() *Prediction
}
