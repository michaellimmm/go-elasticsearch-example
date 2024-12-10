package readingcsv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkReadingCSVWithoutChannel(b *testing.B) {
	b.ResetTimer()
	err := ReadCSVWithoutChannel("sample.csv")
	if err != nil {
		panic(err)
	}
}

func BenchmarkReadingCSVWithChannel(b *testing.B) {
	b.ResetTimer()
	err := ReadCSVWithChannel("sample.csv")
	if err != nil {
		panic(err)
	}
}

func TestReadingCSVWithChannel(t *testing.T) {
	err := ReadCSVWithChannel("sample.csv")
	if err != nil {
		assert.Equal(t, err, nil)
	}
}
