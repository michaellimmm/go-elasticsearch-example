package utils

import (
	"encoding/csv"
	"io"
	"strings"
)

type csvReader struct {
	reader *csv.Reader
	buffer []byte
}

func (r *csvReader) Read(p []byte) (n int, err error) {
	if len(r.buffer) == 0 {
		records, err := r.reader.Read()
		if err != nil {
			return 0, err
		}

		for i, r := range records {
			records[i] = strings.ReplaceAll(r, "\"", "\"\"")
		}

		r.buffer = []byte("\"" + strings.Join(records, "\",\"") + "\"" + "\n")
	}

	n = copy(p, r.buffer)
	r.buffer = r.buffer[n:]
	return n, nil
}

func ConvertCSVReaderToReader(reader *csv.Reader) io.Reader {
	return &csvReader{reader: reader}
}
