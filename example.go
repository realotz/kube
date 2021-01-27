package main

import (
	"fmt"
	"github.com/go-kratos/kube/config"
	"log"
)

func main() {
	s := config.NewSource(config.SourceOption{
		Namespace:     "mesh",
		LabelSelector: "",
	})
	kvs, err := s.Load()
	if err != nil {
		log.Panic(err)
	}
	for _, v := range kvs {
		log.Println(v)
	}
	watcher, err := s.Watch()
	if err != nil {
		log.Panic(err)
	}
	for {
		kvs, err = watcher.Next()
		if err == nil {
			for _, v := range kvs {
				log.Println(v)
			}
		} else {
			fmt.Println(err)
		}
	}
}
