package main

import (
	"fmt"
	"github/shaolim/kakashi/utils/sampler"
)

func main() {
	// Sampling integers from an unknown large stream
	intSampler := sampler.NewReservoirSampler[int](5)

	// Simulate a large stream of integers
	for i := 0; i < 1000; i++ {
		intSampler.Add(i)
	}

	fmt.Println("5 Random Integer Samples:", intSampler.GetSample())

	// Sampling strings from an unknown large stream
	stringNames := []string{
		"Alice", "Bob", "Charlie", "David", "Eve", "Frank",
		"Grace", "Heidi", "Ivan", "Julia", "Kevin", "Linda",
	}

	stringsSampler := sampler.NewReservoirSampler[string](3)

	// Simulate a large stream of names
	for i := 0; i < 100; i++ {
		for _, name := range stringNames {
			stringsSampler.Add(name)
		}
	}

	fmt.Println("3 Random Name Samples:", stringsSampler.GetSample())

	// Custom struct example
	type Person struct {
		Name string
		Age  int
	}

	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
		{Name: "David", Age: 40},
		{Name: "Eve", Age: 28},
	}

	personSampler := sampler.NewReservoirSampler[Person](2)

	// Simulate sampling from a large stream of people
	for i := 0; i < 1000; i++ {
		for _, person := range people {
			personSampler.Add(person)
		}
	}

	fmt.Println("2 Random Person Samples:", personSampler.GetSample())

	datasetSize := 500_000

	fmt.Println("Confidence Interval Sampling:")
	confidenceSampleSize := sampler.CalculateTotalSampleSize(datasetSize, 0.95, 0.05)
	fmt.Printf("Recommended sample size for 95%% confidence, 5%% margin of error: %d\n", confidenceSampleSize)
}
