package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
	_ "github.com/marcboeker/go-duckdb"
)

type item struct {
	LangCode              string  `csv:"Language Code"`
	ID                    string  `csv:"ID"`
	Title                 string  `csv:"Title"`
	Link                  string  `csv:"Link"`
	Price                 float64 `csv:"Price"`
	Currency              string  `csv:"Currency"`
	ImageLink             string  `csv:"Image link"`
	Description           string  `csv:"Description"`
	AdditionalImageLink   string  `csv:"Additional image link"`
	GoogleProductCategory string  `csv:"Google product category"`
	AvailabilityDate      string  `csv:"Availability date"`
	ProductType           string  `csv:"Product type"`
	ProductCode           string  `csv:"Product Code"`
	ProductCodeType       string  `csv:"Product Code type"`
	Condition             string  `csv:"Condition"`
	AgeGroup              string  `csv:"Age group"`
	Color                 string  `csv:"Color"`
	Gender                string  `csv:"Gender"`
	Pattern               string  `csv:"Pattern"`
	Size                  string  `csv:"Size"`
	SizeType              string  `csv:"Size type"`
	SizeSystem            string  `csv:"Size System"`
	Ratings               string  `csv:"Ratings"`
	IsTargetForDelete     string  `csv:"IsTargetForDelete"`
}

type itemScanner struct {
	LangCode              sql.NullString
	ID                    sql.NullString
	Title                 sql.NullString
	Link                  sql.NullString
	Price                 sql.NullString
	Currency              sql.NullString
	Description           sql.NullString
	AdditionalImageLink   sql.NullString
	GoogleProductCategory sql.NullString
	ImageLink             sql.NullString
	AvailabilityDate      sql.NullString
	ProductType           sql.NullString
	ProductCode           sql.NullString
	ProductCodeType       sql.NullString
	Condition             sql.NullString
	AgeGroup              sql.NullString
	Color                 sql.NullString
	Gender                sql.NullString
	Pattern               sql.NullString
	Size                  sql.NullString
	SizeType              sql.NullString
	SizeSystem            sql.NullString
	Ratings               sql.NullString
	IsTargetForDelete     sql.NullString
}

func (i itemScanner) ToItem(idSuffix string) item {
	return item{
		LangCode:              getNullString(i.LangCode),
		ID:                    fmt.Sprintf("%s%s", getNullString(i.ID), idSuffix),
		Title:                 getNullString(i.Title),
		Link:                  getNullString(i.Link),
		Price:                 convStringToFloat64(getNullString(i.Price)),
		Currency:              getNullString(i.Currency),
		Description:           getNullString(i.Description),
		AdditionalImageLink:   getNullString(i.AdditionalImageLink),
		GoogleProductCategory: getNullString(i.GoogleProductCategory),
		ImageLink:             getNullString(i.ImageLink),
		AvailabilityDate:      getNullString(i.AvailabilityDate),
		ProductType:           getNullString(i.ProductType),
		ProductCode:           getNullString(i.ProductCode),
		ProductCodeType:       getNullString(i.ProductCodeType),
		Condition:             getNullString(i.Condition),
		AgeGroup:              getNullString(i.AgeGroup),
		Color:                 getNullString(i.Color),
		Gender:                getNullString(i.Gender),
		Pattern:               getNullString(i.Pattern),
		Size:                  getNullString(i.Size),
		SizeType:              getNullString(i.SizeType),
		SizeSystem:            getNullString(i.SizeSystem),
		Ratings:               getNullString(i.Ratings),
		IsTargetForDelete:     getNullString(i.IsTargetForDelete),
	}
}

func main() {
	totalRow := flag.Int("tr", 0, "Total Rows")
	suffix := flag.String("s", "", "Suffix ID")
	output := flag.String("o", "output.csv", "output filename")
	masterData := flag.String("md", "", "parquet master data ex. masterdata.parquet")

	flag.Parse()

	if *totalRow == 0 || *masterData == "" {
		flag.PrintDefaults()
		return
	}

	db, err := sql.Open("duckdb", "")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := fmt.Sprintf(`select setseed(0.5); SELECT * 
	FROM read_parquet('%s') 
	ORDER BY RANDOM() 
	LIMIT %d`, *masterData, *totalRow)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var items []item

	for rows.Next() {
		var tmp itemScanner

		err := rows.Scan(&tmp.LangCode,
			&tmp.ID,
			&tmp.Title,
			&tmp.Link,
			&tmp.Price,
			&tmp.Currency,
			&tmp.ImageLink,
			&tmp.Description,
			&tmp.AdditionalImageLink,
			&tmp.GoogleProductCategory,
			&tmp.AvailabilityDate,
			&tmp.ProductType,
			&tmp.ProductCode,
			&tmp.ProductCodeType,
			&tmp.Condition,
			&tmp.AgeGroup,
			&tmp.Color,
			&tmp.Gender,
			&tmp.Pattern,
			&tmp.Size,
			&tmp.SizeType,
			&tmp.SizeSystem,
			&tmp.Ratings,
			&tmp.IsTargetForDelete,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		items = append(items, tmp.ToItem(*suffix))
	}

	file, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = gocsv.MarshalFile(&items, file)
	if err != nil {
		log.Fatal(err)
	}
}

func getNullString(ns sql.NullString) string {
	if ns.Valid {
		return strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(ns.String, "\n", "<br>"), "\r", "<br>"))
	}
	return ""
}

func convStringToFloat64(f string) float64 {
	floatValue, _ := strconv.ParseFloat(f, 64)
	return floatValue
}
