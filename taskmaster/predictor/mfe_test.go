package predictor

import (
	"testing"
)

func TestMFE(t *testing.T) {
	// size of 2
	mfe := NewMFE()
	// test that the evict works as expected
	info1 := make(map[string]any)
	info1["fnRequest"] = FnRequest{
		FnName: "a",
		FnParameters: map[string]string{
			"p1": "1",
		},
	}
	info2 := make(map[string]any)
	info2["fnRequest"] = FnRequest{
		FnName: "b",
		FnParameters: map[string]string{
			"p1": "2",
		},
	}
	info3 := make(map[string]any)
	info3["fnRequest"] = FnRequest{
		FnName: "c",
		FnParameters: map[string]string{
			"p1": "3",
		},
	}
	mfe.Update(info1)
	mfe.Update(info2)
	mfe.Update(info3)
	mfe.Update(info2)
	// we expect that mfe currentmaxcount = 2 and currenName should be "b", since we updated it twice
	if mfe.currentMaxCount != 2 {
		t.Errorf("Expected mfe.currentMaxCount to be 2, but it was %d", mfe.currentMaxCount)
	}
	if mfe.currentName != "b" {
		t.Errorf("Expected mfe.currentName to be b, but it was %s", mfe.currentName)
	}
	// now we check that mfe should predict fn b
	prediction := mfe.Predict()
	if prediction.FnName != "b" {
		t.Errorf("Expected fn prediction to be b but it was %s", prediction.FnName)
	}
}
