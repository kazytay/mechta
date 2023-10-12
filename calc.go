package mechta

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
)

type numbers struct {
	A int `json:"a"`
	B int `json:"b"`
}

func Calc(pathToFile string, goroutines int) (int, error) {
	if strings.TrimSpace(pathToFile) == "" {
		return 0, errors.New("pathToFile cannot be empty")
	}
	if goroutines < 1 {
		return 0, errors.New("goroutines cannot be less than 1")
	}

	nums, err := getNumbers(pathToFile)
	if err != nil {
		return 0, err
	}

	var wg sync.WaitGroup
	var sums = make(chan int, goroutines)

	batchSize := len(nums) / goroutines
	for i := 0; i < goroutines; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if i == goroutines-1 {
			end = len(nums)
		}

		wg.Add(1)
		go worker(&wg, nums[start:end], sums)
	}

	go func() {
		wg.Wait()
		close(sums)
	}()

	return getTotal(sums), nil
}

func worker(wg *sync.WaitGroup, nums []numbers, sums chan int) {
	defer wg.Done()

	var sum int
	for _, v := range nums {
		sum += v.A + v.B
	}
	sums <- sum
}

func getTotal(sums chan int) (total int) {
	for sum := range sums {
		total += sum
	}
	return
}

func getNumbers(path string) (nums []numbers, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&nums)
	if err != nil {
		return
	}
	return
}
