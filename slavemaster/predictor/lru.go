package predictor

import (
	"container/list"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type LRUConfig struct {
	Size int64 `yaml:"size"`
}

// LRU is a Predictor that predicts that the least recently used function will be invoked again
// internally, it stores a list
// if say the functions A, B, C are invoked
// the list will be front -> C -> B -> A
// when it is polled for a prediction, it will predict that A (least recently used) should be pinged.
// this makes sense because A will likely have expired to timeout
// say it has a max size (n) of 3
// then if D is invoked, lru should remove C
// front -> D -> B -> A
// so the eviction really is MRU because the function container of most benefit is the LRU
type LRU struct {
	// max size
	size int64
	// current size
	n   int64
	lst *list.List
	// map of functions to their fnParameters
	fnSet map[string]map[string]string
}

func NewLRU(lru_config_file string) *LRU {
	yaml_file, err := os.ReadFile(lru_config_file)
	if err != nil {
		log.Fatal(err)
	}
	var config LRUConfig
	err = yaml.Unmarshal(yaml_file, &config)
	if err != nil {
		log.Fatal(err)
	}
	var lru = new(LRU)
	lru.size = config.Size
	lru.n = 0
	lru.lst = list.New()
	lru.fnSet = make(map[string]map[string]string)
	return lru
}

// add adds the fnRequest to the linked list, evicting if necessary
func (lru *LRU) add(request FnRequest) {
	// maximum size of the list, evict
	// we evict the first item in the list
	if lru.n == lru.size {
		// remove from the list
		toRemove := lru.lst.Front()
		lru.lst.Remove(toRemove)
		lru.n--
		// now remove from the fnSet
		delete(lru.fnSet, request.FnName)
	}
	// add to front of list
	lru.lst.PushFront(request.FnName)
	lru.n++
	// update fnParameters
	lru.fnSet[request.FnName] = request.FnParameters
}

func (lru *LRU) Update(info map[string]any) {
	// map should contain
	//  1. fnRequest: FnRequest struct
	var fnRequest FnRequest
	fnRequest = info["fnRequest"].(FnRequest)
	//:= info["fnRequest"]
	// LRU algorithm
	// First check if it exists in the list
	_, contains := lru.fnSet[fnRequest.FnName]
	if !contains { // if it does not exist in the list, simple add
		lru.add(fnRequest)
	} else { // if it exists in the list, move it to the front of the list
		var e *list.Element
		for e = lru.lst.Front(); e.Value != fnRequest.FnName; e = e.Next() {
		}
		lru.lst.MoveToFront(e)
	}
	// finally, update the fnParameters
	lru.fnSet[fnRequest.FnName] = fnRequest.FnParameters
}

// Predict returns the last item in the internal linked list
func (lru *LRU) Predict() *Prediction {
	if lru.n == 0 {
		return NilPrediction
	}
	e := lru.lst.Back()
	fnName := e.Value.(string)
	return &Prediction{
		FnName:       fnName,
		FnParameters: lru.fnSet[fnName],
	}
}
