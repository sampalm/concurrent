package main

import (
	// "concurrency/throttle"

	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

var sema = semaphore.NewWeighted(5)

func main() {
	path := os.Getenv("LOCAL_PATH")
	chanSize := make(chan int64)

	th := &sync.WaitGroup{}
	t := time.Now()

	th.Add(1)
	go ListDir(path, th, chanSize)

	go func() {
		th.Wait()
		close(chanSize)

	}()

	var totalSize int64
	var totalFiles int64

	for size := range chanSize {
		totalSize += size
		totalFiles++
	}

	duration := time.Since(t)

	fmt.Println("COUNT FILES ===> ", totalFiles)
	fmt.Println("COUNT SIZE ===> ", totalSize)
	fmt.Println("DURATION ===> ", duration)
}

func ListDir(path string, th *sync.WaitGroup, chSize chan<- int64) {
	var size int64

	if err := sema.Acquire(context.Background(), 1); err != nil {
		log.Fatal(err)
	}
	defer func() { sema.Release(1) }()

	defer th.Done()

	entry, err := os.ReadDir(path)
	if err != nil {
		return
	}

	for _, subEntry := range entry {
		if subEntry.IsDir() {
			th.Add(1)
			go ListDir(filepath.Join(path, subEntry.Name()), th, chSize)
		} else {
			f, _ := subEntry.Info()
			size += f.Size()
			chSize <- size
		}
	}

}
