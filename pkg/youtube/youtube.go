package youtube

import (
	"fmt"
	"sync"
)

var lock = &sync.Mutex{}

type youtube struct {
}

var client *youtube

func GetClient() *youtube {
	if client == nil {
		lock.Lock()
		defer lock.Unlock()
		if client == nil {
			fmt.Println("Creating youtube client now.")
			client = &youtube{}
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return client
}

func Search() {
	fmt.Println("searching youtube")
	GetClient()
}
