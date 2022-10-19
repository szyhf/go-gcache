package main

import (
	"fmt"

	"github.com/szyhf/go-gcache/v2"
)

func main() {
	gc := gcache.New[string, string](10).
		LFU().
		LoaderFunc(func(key string) (string, error) {
			return fmt.Sprintf("%v-value", key), nil
		}).
		Build()

	v, err := gc.Get("key")
	if err != nil {
		panic(err)
	}
	fmt.Println(v)
}
