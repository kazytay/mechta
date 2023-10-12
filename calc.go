package mechta

import (
	"encoding/json"
	"os"
	"sync"
)

type numbers struct {
	A int `json:"a"`
	B int `json:"b"`
}

func Calc(goroutines int) (*int, error) {
	nums, err := getNumbers()
	if err != nil {
		return nil, err
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

	total := getTotal(sums)
	return &total, nil
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

func getNumbers() (nums []numbers, err error) {
	f, err := os.Open("numbers.json")
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
