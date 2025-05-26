package main

import (
	"log"
	"time"

	"github.com/szyhf/go-gcache/v2"
)

func main() {
	log.SetFlags(log.Lshortfile)
	cache := gcache.New[int, int](100).
		PurgeVisitorFunc(func(k int, v int) {
			log.Println(k, "=", v)
		}).
		LoaderExpireFunc(func(k int) (int, *time.Duration, error) {
			log.Printf("load(k=%d)", k)
			t := time.Second
			return k, &t, nil
		}).
		EvictedFunc(func(k int, v int) {
			log.Printf("回收(k=%d,v=%d)", k, v)
		}).
		Lease(time.Second).
		ARC().
		Build()
	go func() {
		for {
			time.Sleep(time.Second / 2)
			cache.Get(99)
		}
	}()
	for i := 0; i < 100; i++ {
		cache.SetWithExpire(i, i, time.Second)
		// cache.SetWithExpire(i, i, 0) //time.Second)
	}
	log.Println(cache.Get(99))
	time.Sleep(time.Second * 2)
	log.Println(cache.Get(99))
	log.Println(cache.Get(99))
}
