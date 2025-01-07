package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

type Item struct {
	LangCode              string `parquet:"name=langcode, type=BYTE_ARRAY, convertedtype=UTF8"`
	ID                    string `parquet:"name=id, type=BYTE_ARRAY, convertedtype=UTF8"`
	Title                 string `parquet:"name=title, type=BYTE_ARRAY, convertedtype=UTF8"`
	Link                  string `parquet:"name=link, type=BYTE_ARRAY, convertedtype=UTF8"`
	Price                 string `parquet:"name=price, type=BYTE_ARRAY, convertedtype=UTF8"`
	Currency              string `parquet:"name=currency, type=BYTE_ARRAY, convertedtype=UTF8"`
	ImageLink             string `parquet:"name=imagelink, type=BYTE_ARRAY, convertedtype=UTF8"`
	Description           string `parquet:"name=description, type=BYTE_ARRAY, convertedtype=UTF8"`
	AdditionalImageLink   string `parquet:"name=additionalimagelnk, type=BYTE_ARRAY, convertedtype=UTF8"`
	GoogleProductCategory string `parquet:"name=googleproductcategory, type=BYTE_ARRAY, convertedtype=UTF8"`
	AvailabilityDate      string `parquet:"name=availabilitydate, type=BYTE_ARRAY, convertedtype=UTF8"`
	ProductType           string `parquet:"name=producttype, type=BYTE_ARRAY, convertedtype=UTF8"`
	ProductCode           string `parquet:"name=productcode, type=BYTE_ARRAY, convertedtype=UTF8"`
	ProductCodeType       string `parquet:"name=productcodetype, type=BYTE_ARRAY, convertedtype=UTF8"`
	Condition             string `parquet:"name=condition, type=BYTE_ARRAY, convertedtype=UTF8"`
	AgeGroup              string `parquet:"name=agegroup, type=BYTE_ARRAY, convertedtype=UTF8"`
	Color                 string `parquet:"name=color, type=BYTE_ARRAY, convertedtype=UTF8"`
	Gender                string `parquet:"name=gender, type=BYTE_ARRAY, convertedtype=UTF8"`
	Pattern               string `parquet:"name=pattern, type=BYTE_ARRAY, convertedtype=UTF8"`
	Size                  string `parquet:"name=size, type=BYTE_ARRAY, convertedtype=UTF8"`
	SizeType              string `parquet:"name=sizetype, type=BYTE_ARRAY, convertedtype=UTF8"`
	SizeSystem            string `parquet:"name=sizesystem, type=BYTE_ARRAY, convertedtype=UTF8"`
	Ratings               string `parquet:"name=ratings, type=BYTE_ARRAY, convertedtype=UTF8"`
	IsTargetForDelete     string `parquet:"name=istargetfordelete, type=BYTE_ARRAY, convertedtype=UTF8"`
}

func main() {
	inputFile := flag.String("input", "", "input CSV file")
	outputFile := flag.String("output", "", "output parquet file")

	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		flag.PrintDefaults()
		return
	}

	start := time.Now()

	csvFile, err := os.Open(*inputFile)
	if err != nil {
		log.Fatal("Error opening CSV file:", err)
	}
	defer csvFile.Close()

	fileInfo, err := csvFile.Stat()
	if err != nil {
		log.Fatal("Error getting file info:", err)
	}
	totalSize := fileInfo.Size()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read CSV header
	_, err = reader.Read()
	if err != nil {
		log.Fatal("Error reading CSV header:", err)
	}

	// Create Parquet file
	fw, err := os.Create(*outputFile)
	if err != nil {
		log.Fatal("Error creating Parquet file:", err)
	}
	defer fw.Close()

	// Configure Parquet writer
	numCPU := runtime.NumCPU()
	pw, err := writer.NewParquetWriterFromWriter(fw, new(Item), int64(numCPU))
	if err != nil {
		log.Fatal("Error creating Parquet writer:", err)
	}

	// Set optimized settings for large files
	pw.CompressionType = parquet.CompressionCodec_ZSTD // Use ZSTD compression
	pw.RowGroupSize = 128 * 1024 * 1024                // 128MB row groups
	pw.PageSize = 1 * 1024 * 1024                      // 1MB pages

	// Read and convert records
	lineCount := 0
	reportInterval := 100000 // Report progress every 100k records
	lastReportTime := time.Now()

	for {
		row, err := reader.Read()
		if err != nil {
			break // End of file
		}
		lineCount++

		// Progress reporting
		if lineCount%reportInterval == 0 {
			currentTime := time.Now()
			duration := currentTime.Sub(lastReportTime)
			recordsPerSecond := float64(reportInterval) / duration.Seconds()
			fmt.Printf("Processed %d records (%.2f records/sec)\n", lineCount, recordsPerSecond)
			lastReportTime = currentTime
		}

		// Create record
		item := Item{
			LangCode:              row[0],
			ID:                    row[1],
			Title:                 row[2],
			Link:                  row[3],
			Price:                 row[4],
			Currency:              row[5],
			ImageLink:             row[6],
			Description:           row[7],
			AdditionalImageLink:   row[8],
			GoogleProductCategory: row[9],
			AvailabilityDate:      row[10],
			ProductType:           row[11],
			ProductCode:           row[12],
			ProductCodeType:       row[13],
			Condition:             row[14],
			AgeGroup:              row[15],
			Color:                 row[16],
			Gender:                row[17],
			Pattern:               row[18],
			Size:                  row[19],
			SizeType:              row[20],
			SizeSystem:            row[21],
			Ratings:               row[22],
			IsTargetForDelete:     row[23],
		}

		if err := pw.Write(item); err != nil {
			log.Printf("Error writing record on line %d: %v", lineCount, err)
			continue
		}
	}

	// Close writer
	if err := pw.WriteStop(); err != nil {
		log.Fatal("Error closing writer:", err)
	}

	// Get final file sizes
	parquetInfo, err := os.Stat(*outputFile)
	if err != nil {
		log.Fatal("Error getting parquet file info:", err)
	}

	// Print summary
	duration := time.Since(start)
	compressionRatio := float64(totalSize) / float64(parquetInfo.Size())

	fmt.Printf("\nConversion Summary:\n")
	fmt.Printf("Total records processed: %d\n", lineCount)
	fmt.Printf("Input file size: %.2f MB\n", float64(totalSize)/(1024*1024))
	fmt.Printf("Output file size: %.2f MB\n", float64(parquetInfo.Size())/(1024*1024))
	fmt.Printf("Compression ratio: %.2fx\n", compressionRatio)
	fmt.Printf("Processing time: %v\n", duration)
	fmt.Printf("Average speed: %.2f records/second\n", float64(lineCount)/duration.Seconds())
}
