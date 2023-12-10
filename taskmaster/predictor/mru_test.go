package predictor

import (
	"testing"
)

func TestMRU(t *testing.T) {
	// size of 2
	mru := NewMRU("../mru_config.yaml")
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
	mru.Update(info1)
	mru.Update(info2)
	mru.Update(info3)
	prediction := mru.Predict()
	if prediction.FnName != "c" {
		t.Errorf("Expected fn prediction to be c but it was %s", prediction.FnName)
	}
}
