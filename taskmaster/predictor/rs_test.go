package predictor

import (
	"testing"
)

func TestRS(t *testing.T) {
	// size of 2
	rs := NewRS()
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
	rs.Update(info1)
	rs.Update(info2)
	rs.Update(info3)
	rs.Update(info2)

	if len(rs.functions) != 3 {
		t.Errorf("Expected len(rs.functions) to be 3, but it was %d", len(rs.functions))
	}
	prediction := rs.Predict()
	print(prediction.FnName)
}
