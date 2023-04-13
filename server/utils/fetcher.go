package utils

import (
	"time"
)

type cacheUnit[T any] struct {
	value *T
	ts    time.Time
}

type Fetcher[KeyType comparable, ValueType any] struct {
	cache map[KeyType]cacheUnit[ValueType]
	waits map[KeyType][]*Future[ValueType]

	cacheInterval time.Duration

	jobs  chan func()
	fetch func(KeyType) (*ValueType, error)
}

func NewFetcher[KeyType comparable, ValueType any](interval time.Duration, fetch func(KeyType) (*ValueType, error)) *Fetcher[KeyType, ValueType] {
	ret := &Fetcher[KeyType, ValueType]{
		cache:         make(map[KeyType]cacheUnit[ValueType]),
		waits:         make(map[KeyType][]*Future[ValueType]),
		cacheInterval: interval,
		jobs:          make(chan func(), 128),
		fetch:         fetch,
	}
	go ret.run()
	return ret
}

func (f *Fetcher[KeyType, ValueType]) run() {
	for job := range f.jobs {
		job()
	}
}

func (f *Fetcher[KeyType, ValueType]) postJob(job func()) {
	f.jobs <- job
}

func (f *Fetcher[KeyType, ValueType]) Get(key KeyType) (*ValueType, error) {
	future := MakeFuture[ValueType]()
	f.postJob(func() {
		unit, ok := f.cache[key]
		if ok && time.Since(unit.ts) < f.cacheInterval {
			future.Done(unit.value, nil)
			return
		}

		f.waits[key] = append(f.waits[key], future)
		if len(f.waits[key]) == 1 {
			go func() {
				ret, err := f.fetch(key)
				f.postJob(func() {
					if err == nil {
						f.cache[key] = cacheUnit[ValueType]{
							value: ret,
							ts:    time.Now(),
						}
					}
					futures := f.waits[key]
					f.waits[key] = []*Future[ValueType]{}
					for _, item := range futures {
						item.Done(ret, err)
					}
				})
			}()
		}
	})
	return future.Wait()
}
