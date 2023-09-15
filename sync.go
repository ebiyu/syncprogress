package main

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
)

func main() {
	// start progress bar
	completed := make(chan bool)
	go showDiskCacheProgress(completed)

	// sync
	output, err := exec.Command("sh", "-c", "sync").CombinedOutput()

	completed <- true
	time.Sleep(time.Second / 2 * 3)

	if err != nil {
		log.Println(string(output))
		panic(err)
	}
}

func getRemainDiskCache() (int, error) {
	result, err := exec.Command("sh", "-c", " grep -e 'Dirty' /proc/meminfo | awk '{print $2}'").Output()
	if err != nil {
		return 0, err
	}

	kb, err := strconv.Atoi(strings.Trim(string(result), "\n"))
	if err != nil {
		return 0, err
	}

	return kb, nil
}

func showDiskCacheProgress(completed chan bool) {
	first, err := getRemainDiskCache()
	if err != nil {
		log.Fatal(err)
	}

	bar := pb.StartNew(first * 1000)
	bar.Set(pb.Bytes, true)

	time.Sleep(time.Second / 2)

	for {
		kb, err := getRemainDiskCache()
		if err != nil {
			continue
		}

		bar.SetCurrent(int64(first-kb) * 1000)

		select {
		case i := <-completed:
			if i {
				bar.Finish()
				return
			}

		default:
			//fmt.Println("No value")
		}
	}
}
