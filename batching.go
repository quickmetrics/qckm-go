package qm

import (
	"log"
	"sync"
	"time"
)

// batch
type batch struct {
	mutex     *sync.RWMutex
	maxSize   int
	maxWait   time.Duration
	items     []event
	flushChan chan []event
	wg        *sync.WaitGroup
}

// initialize batcher
func newBatcher(maxSize int, maxWait time.Duration, workerSize int) *batch {
	// make sure our settings are sensible
	size, wait, workers := cleanSettings(maxSize, maxWait, workerSize)

	instance := &batch{
		maxSize:   size,
		maxWait:   wait,
		items:     make([]event, 0),
		mutex:     &sync.RWMutex{},
		flushChan: make(chan []event, workers),
		wg:        &sync.WaitGroup{},
	}
	instance.setFlushWorker(workers)
	go instance.runFlushByTime()
	return instance
}

func cleanSettings(maxSize int, maxWait time.Duration, workerSize int) (int, time.Duration, int) {
	if maxSize < 10 {
		maxSize = 10
	} else if maxSize > 10000 {
		maxSize = 10000
	}

	if maxWait < 1 {
		maxWait = 1
	} else if maxWait > 60 {
		maxWait = 60
	}
	maxWait = maxWait * time.Second

	if workerSize < 1 {
		workerSize = 1
	} else if workerSize > 10 {
		workerSize = 10
	}

	if isVerbose {
		log.Printf("init batcher with maxSize: %v, maxWait: %v, workers: %v", maxSize, maxWait, workerSize)
	}

	return maxSize, maxWait, workerSize
}

func (b *batch) process(ee []event) {
	processBatch(ee, b.wg)
}

func (b *batch) setFlushWorker(workerSize int) {
	if workerSize < 1 {
		workerSize = 1
	}
	for id := 1; id <= workerSize; id++ {
		go func(workerID int, flushJobs <-chan []event) {
			for j := range flushJobs {
				b.process(j)
			}
		}(id, b.flushChan)
	}
}

func (b *batch) add(data event) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.items = append(b.items, data)

	b.wg.Add(1)

	if len(b.items) >= b.maxSize {
		b.flush()
	}
}

func (b *batch) runFlushByTime() {
	for {
		select {
		case <-time.After(b.maxWait):
			b.mutex.Lock()
			b.flush()
			b.mutex.Unlock()
		}
	}
}

func (b *batch) wait() {
	b.wg.Wait()
}

func (b *batch) flush() {
	if len(b.items) <= 0 {
		return
	}

	copiedItems := make([]event, len(b.items))
	for idx, i := range b.items {
		copiedItems[idx] = i
	}
	b.items = b.items[:0]
	b.flushChan <- copiedItems
}
