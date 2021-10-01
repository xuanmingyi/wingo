package register

import (
	"fmt"
	"sync"
	"testing"
)

func RecordCount(r *Register, wg *sync.WaitGroup) {
	for i := 0; i < 1000; i++ {
		record, err := r.Get("count")
		if err != nil {
			panic(err)
		}

		count := record.(int)

		r.Set("count", count+1)
	}
	wg.Done()
}

func TestRegister(t *testing.T) {
	r := Register{}

	r.Open("sqlite://a.db")

	defer r.Close()

	var record Record
	var err error
	err = r.Set("ssss", 13)
	if err != nil {
		panic(err)
	}

	err = r.Set("ssss", 15)
	if err != nil {
		panic(err)
	}

	record, err = r.Get("ssss")
	if err != nil {
		panic(err)
	}
	fmt.Println(record)

	/*
		r.Set("count", 0)
		var wg sync.WaitGroup

		wg.Add(5)

		go RecordCount(&r, &wg)
		go RecordCount(&r, &wg)
		go RecordCount(&r, &wg)
		go RecordCount(&r, &wg)
		go RecordCount(&r, &wg)

		wg.Wait()
	*/
}
