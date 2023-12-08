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
	// // we expect that lru should only have functions a and c, cos b is MRU and should be ignored
	// ll := mru.lst
	// var e *list.Element
	// for e = ll.Front(); e != nil; e = e.Next() {
	// 	fmt.Println("element: ", e.Value)
	// 	if e.Value == "b" {
	// 		t.Error("Expect c to not be in lru but encountered it")
	// 	}
	// }
	// now we check that lru should predict fn a
	prediction := mru.Predict()
	if prediction.FnName != "c" {
		t.Errorf("Expected fn prediction to be c but it was %s", prediction.FnName)
	}
}
