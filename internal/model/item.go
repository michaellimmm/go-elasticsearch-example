package model

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	LanguageCode          string  `csv:"Language Code"`
	Id                    string  `csv:"ID"`
	Title                 string  `csv:"Title"`
	Link                  string  `csv:"Link"`
	Price                 string  `csv:"Price"`
	CurrencyCode          string  `csv:"Currency"`
	ImageLink             string  `csv:"Image link"`
	Description           string  `csv:"Description"`
	AdditionalImageLink   string  `csv:"Additional image link"`
	GoogleProductCategory string  `csv:"Google product category"`
	AvailableFrom         string  `csv:"Availability date"`
	ProductType           string  `csv:"Product type"`
	ProductCode           string  `csv:"Product Code"`
	ProductCodeType       string  `csv:"Product Code type"`
	Condition             string  `csv:"Condition"`
	AgeGroup              string  `csv:"Age group"`
	Color                 string  `csv:"Color"`
	Gender                string  `csv:"Gender"`
	Pattern               string  `csv:"Pattern"`
	SizeValue             string  `csv:"Size"`
	SizeType              string  `csv:"Size type"`
	SizeSystem            string  `csv:"Size system"`
	Ratings               float64 `csv:"Ratings"`
	IsTargetForDelete     string  `csv:"IsTargetForDelete"`
}

func (i Item) IsDeleted() bool {
	return i.IsTargetForDelete == "1"
}

type ItemDoc struct {
	LanguageCode         string                 `json:"languageCode"`
	ID                   primitive.ObjectID     `json:"mongoId,omitempty"`
	Sku                  string                 `json:"sku"`
	Title                string                 `json:"title"`
	Link                 string                 `json:"link"`
	Price                *Price                 `json:"price"`
	Images               []string               `json:"images"`
	Description          string                 `json:"description"`
	IsDeleted            bool                   `json:"isDeleted"`
	Record               *RecordWithDelete      `json:"record"`
	AdditionalProperties map[string]interface{} `json:"additionalProperties"`
}

type Price struct {
	CurrencyCode string `json:"currencyCode"`
	PriceMajor   uint32 `json:"priceMajor"`
	PriceMinor   uint32 `json:"priceMinor"`
}

type RecordWithDelete struct {
	Created time.Time
	Updated time.Time
	Deleted *time.Time
}

func ConvertItemToItemDoc(item Item) ItemDoc {
	var (
		priceMajor uint64
		priceMinor uint64
	)

	splitPrice := strings.Split(item.Price, ".")
	if len(splitPrice) > 1 {
		priceMajor, _ = strconv.ParseUint(splitPrice[0], 10, 32)
		priceMinor, _ = strconv.ParseUint(splitPrice[1], 10, 32)
	} else if len(splitPrice) > 0 {
		priceMajor, _ = strconv.ParseUint(splitPrice[0], 10, 32)
	}

	images := []string{}
	if item.ImageLink != "" {
		images = append(images, item.ImageLink)
	}
	if item.AdditionalImageLink != "" {
		for _, imageUrl := range strings.Split(item.AdditionalImageLink, " ") {
			_, err := url.ParseRequestURI(strings.TrimSpace(imageUrl))
			if err != nil {
				continue
			}
			images = append(images, imageUrl)
		}
	}

	return ItemDoc{
		LanguageCode: item.LanguageCode,
		ID:           primitive.NewObjectID(),
		Sku:          item.Id,
		Title:        item.Title,
		Link:         item.Link,
		Price: &Price{
			CurrencyCode: item.CurrencyCode,
			PriceMajor:   uint32(priceMajor),
			PriceMinor:   uint32(priceMinor),
		},
		Images:      images,
		Description: item.Description,
		IsDeleted:   item.IsDeleted(),
		Record:      &RecordWithDelete{Created: time.Now(), Updated: time.Now()},
		AdditionalProperties: map[string]interface{}{
			"GoogleProductCategory": item.GoogleProductCategory,
			"AvailableFrom":         item.AvailableFrom,
			"ProductType":           item.ProductType,
			"Condition":             item.Condition,
			"AgeGroup":              item.AgeGroup,
			"Color":                 item.Color,
			"Gender":                item.Gender,
			"Pattern":               item.Pattern,
			"SizeValue":             item.SizeValue,
			"SizeType":              item.SizeType,
			"SizeSystem":            item.SizeSystem,
			"Ratings":               item.Ratings,
			"ProductCodeType":       item.ProductCodeType,
			"ProductCode":           item.ProductCode,
			"Title":                 item.Title,
			"Description":           item.Description,
			"CurrencyCode":          item.CurrencyCode,
			"PriceMajor":            uint32(priceMajor),
			"PriceMinor":            uint32(priceMinor),
		},
	}
}
