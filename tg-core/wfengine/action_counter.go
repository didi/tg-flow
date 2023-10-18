package wfengine

import "sync"

type ActionCounter struct {
	wgMap 					map[string]*sync.WaitGroup
	skipedActionIdPairs		*sync.Map

}
