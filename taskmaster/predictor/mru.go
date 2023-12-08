package predictor

import (
	"container/list"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type MRUConfig struct {
	Size int64 `yaml:"size"`
}

// mru is a Predictor that predicts that the least recently used function will be invoked again
// internally, it stores a list
// if say the functions A, B, C are invoked
// the list will be front -> C -> B -> A
// when it is polled for a prediction, it will predict that A (least recently used) should be pinged.
// this makes sense because A will likely have expired to timeout
// say it has a max size (n) of 3
// then if D is invoked, mru should remove C
// front -> D -> B -> A
// so the eviction really is MRU because the function container of most benefit is the mru
type MRU struct {
	// max size
	size int64
	// current size
	n   int64
	lst *list.List
	// map of functions to their fnParameters
	fnSet map[string]map[string]string
}

func NewMRU(mru_config_file string) *MRU {
	yaml_file, err := os.ReadFile(mru_config_file)
	if err != nil {
		log.Fatal(err)
	}
	var config MRUConfig
	err = yaml.Unmarshal(yaml_file, &config)
	if err != nil {
		log.Fatal(err)
	}
	var mru = new(MRU)
	mru.size = config.Size
	mru.n = 0
	mru.lst = list.New()
	mru.fnSet = make(map[string]map[string]string)
	return mru
}

// add adds the fnRequest to the linked list, evicting if necessary
func (mru *MRU) add(request FnRequest) {
	// maximum size of the list, evict
	// we evict the first item in the list
	if mru.n == mru.size {
		// remove from the list
		toRemove := mru.lst.Front()
		mru.lst.Remove(toRemove)
		mru.n--
		// now remove from the fnSet
		delete(mru.fnSet, request.FnName)
	}
	// add to front of list
	mru.lst.PushFront(request.FnName)
	mru.n++
	// update fnParameters
	mru.fnSet[request.FnName] = request.FnParameters
}

func (mru *MRU) Update(info map[string]any) {
	// map should contain
	//  1. fnRequest: FnRequest struct
	var fnRequest FnRequest = info["fnRequest"].(FnRequest)
	log.Printf("%+v\n", fnRequest)
	// mru algorithm
	// First check if it exists in the list
	_, contains := mru.fnSet[fnRequest.FnName]
	log.Printf("Contains = %v\n", contains)
	if !contains { // if it does not exist in the list, simple add
		mru.add(fnRequest)
	} else { // if it exists in the list, move it to the front of the list
		var e *list.Element
		// FIXME: Investigate why e can still be nil?
		// maybe the set is not being updated correctly
		for e = mru.lst.Front(); e != nil && e.Value != fnRequest.FnName; e = e.Next() {
		}
		// Somehow, it's possible that e cud be nil?
		if e == nil {
			mru.add(fnRequest)
		} else {
			mru.lst.MoveToFront(e)
		}
	}
	// finally, update the fnParameters
	mru.fnSet[fnRequest.FnName] = fnRequest.FnParameters
}

// Predict returns the last item in the internal linked list
func (mru *MRU) Predict() *Prediction {
	if mru.n == 0 {
		return NilPrediction
	}
	e := mru.lst.Front()
	fnName := e.Value.(string)
	return &Prediction{
		FnName:       fnName,
		FnParameters: mru.fnSet[fnName],
	}
}
