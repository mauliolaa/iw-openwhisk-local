package predictor

import (
	"testing"
)

func TestPQueue(t *testing.T) {
	// size of 2
	pq := make(PriorityQueue, 0)
	// test that the evict works as expected
	info1 := make(map[string]any)
	info1["fnRequest"] = FnRequest{
		FnName: "a",
		FnParameters: map[string]string{
			"p1": "1",
		},
	}
	info1["language"] = "js"
	info2 := make(map[string]any)
	info2["fnRequest"] = FnRequest{
		FnName: "b",
		FnParameters: map[string]string{
			"p1": "2",
		},
	}
	info2["language"] = "go"
	info3 := make(map[string]any)
	info3["fnRequest"] = FnRequest{
		FnName: "c",
		FnParameters: map[string]string{
			"p1": "3",
		},
	}
	info3["language"] = "jar"
	info4 := make(map[string]any)
	info4["fnRequest"] = FnRequest{
		FnName: "d",
		FnParameters: map[string]string{
			"p1": "3",
		},
	}
	info4["language"] = "rb"
	pq.Update(info1)
	pq.Update(info2)
	pq.Update(info3)
	pq.Update(info4)
	print(pq.Len())
	// info3 should be popped first since java has the highest penalty/priority
	// now we check that lru should predict fn a
	prediction := pq.Predict()
	if prediction.FnName != "c" {
		t.Errorf("Expected fn prediction to be c but it was %s", prediction.FnName)
	}
	// info4 should go next as ruby has second highest penalty
	prediction = pq.Predict()
	if prediction.FnName != "d" {
		t.Errorf("Expected fn prediction to be a but it was %s", prediction.FnName)
	}
	// info1 should go next as it has same priority as info2 but came first
	prediction = pq.Predict()
	if prediction.FnName != "a" {
		t.Errorf("Expected fn prediction to be a but it was %s", prediction.FnName)
	}
	// now info2
	prediction = pq.Predict()
	if prediction.FnName != "b" {
		t.Errorf("Expected fn prediction to be b but it was %s", prediction.FnName)
	}
	// Now should be empty
	if pq.Len() != 0 {
		t.Errorf("Expected pq to be empty but it was not!")
	}
}
