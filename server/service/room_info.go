package service

import (
	"time"

	dm "github.com/zerozwt/BLiveDanmaku"
	"github.com/zerozwt/blivehl/server/logger"
)

var cacheInterval time.Duration = time.Second * 3

type roomInfoCache struct {
	info *dm.RoomInfo
	ts   time.Time
}

type fetcherFuture struct {
	wait_ch chan bool
	ret     *dm.RoomInfo
	err     error
}

func newFetcherFuture() *fetcherFuture {
	return &fetcherFuture{
		wait_ch: make(chan bool),
		ret:     nil,
		err:     nil,
	}
}

func (f *fetcherFuture) Wait() (*dm.RoomInfo, error) {
	<-f.wait_ch
	return f.ret, f.err
}

func (f *fetcherFuture) Done(ret *dm.RoomInfo, err error) {
	f.ret = ret
	f.err = err
	close(f.wait_ch)
}

type RoomInfoFetcher struct {
	jobs  chan func()
	cache map[int]*roomInfoCache
	queue map[int][]*fetcherFuture
}

var gFetcher *RoomInfoFetcher = newRoomInfoFetcher()

func newRoomInfoFetcher() *RoomInfoFetcher {
	ret := &RoomInfoFetcher{
		jobs:  make(chan func(), 1024),
		cache: make(map[int]*roomInfoCache),
		queue: make(map[int][]*fetcherFuture),
	}
	go ret.run()
	return ret
}

func GetRoomInfoFetcher() *RoomInfoFetcher {
	return gFetcher
}

func (f *RoomInfoFetcher) run() {
	for job := range f.jobs {
		job()
	}
}

func (f *RoomInfoFetcher) postAndWait(job func()) {
	wait_ch := make(chan bool)
	wrap := func() {
		defer func() {
			close(wait_ch)
			if err := recover(); err != nil {
				logger.ERROR("RoomInfoFetcher run job panic: %v", err)
			}
		}()
		job()
	}
	f.jobs <- wrap
	<-wait_ch
}

func (f *RoomInfoFetcher) GetRoomInfo(roomId int) (*dm.RoomInfo, error) {
	future := newFetcherFuture()

	f.postAndWait(func() {
		// check if cached data is available
		unit, ok := f.cache[roomId]
		if ok && time.Since(unit.ts) < cacheInterval {
			future.Done(unit.info, nil)
			return
		}

		// put future into the queue
		f.queue[roomId] = append(f.queue[roomId], future)
		if len(f.queue[roomId]) > 1 {
			return
		}

		// start a request
		go func() {
			info, err := dm.GetRoomInfo(roomId)
			f.postAndWait(func() {
				if err == nil {
					f.cache[roomId] = &roomInfoCache{
						info: info,
						ts:   time.Now(),
					}
				}
				futures := f.queue[roomId]
				f.queue[roomId] = []*fetcherFuture{}
				for _, future := range futures {
					future.Done(info, err)
				}
			})
		}()
	})

	return future.Wait()
}
