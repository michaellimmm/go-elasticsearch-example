package readingcsv

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"

	"github.com/gocarina/gocsv"

	"github/shaolim/go-elasticsearch-example/app/cli/model"
)

func ReadCSVWithoutChannel(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var items []*model.Item
	if err := gocsv.UnmarshalCSV(reader, &items); err != nil {
		fmt.Printf("failed to UnmarshalCSV: %+v\n", err)
		return err
	}

	count := 0
	for range items {
		count++
	}

	fmt.Printf("without channel count: %d\n", count)

	return nil
}

func ReadCSVWithChannel(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		return err
	}
	defer file.Close()

	queue := make(chan *model.Item, 1000)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		processItem(queue)
	}()

	go func() {
		if err := gocsv.UnmarshalToChan(file, queue); err != nil {
			fmt.Printf("failed to UnmarshalToChan: %+v\n", err)
			return
		}
	}()

	wg.Wait()

	return nil
}

func processItem(in <-chan *model.Item) {
	count := 0
	for range in {
		count++
	}
	fmt.Printf("with channel count: %d\n", count)
}
