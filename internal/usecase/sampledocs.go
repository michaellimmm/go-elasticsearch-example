package usecase

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/gocarina/gocsv"

	"github/shaolim/kakashi/internal/model"
	"github/shaolim/kakashi/pkg/esclient"
	"github/shaolim/kakashi/pkg/esclient/esquery"
	"github/shaolim/kakashi/utils/sampler"
)

type SampleDocs struct {
	esclient esclient.Client
}

func NewSampleDocsUseCase(esClient esclient.Client) *SampleDocs {
	return &SampleDocs{
		esclient: esClient,
	}
}

func (s *SampleDocs) Execute(index string, filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var items []*model.Item
	if err := gocsv.UnmarshalCSV(reader, &items); err != nil {
		fmt.Printf("failed to UnmarshalCSV: %+v\n", err)
		return err
	}

	totalRows := len(items)
	totalSample := sampler.CalculateTotalSampleSize(totalRows, 0.95, 0.05)

	rs := sampler.NewReservoirSampler[*model.Item](totalSample)
	for _, item := range items {
		rs.Add(item)
	}

	sample := rs.GetSample()
	fmt.Printf("total rows: %d, sample size: %d\n", totalRows, len(sample))

	termQueries := make([]esquery.QueryType, 0, len(sample))
	for _, item := range sample {
		termQueries = append(termQueries, esquery.Term("_id", item.Id))
	}

	res, err := s.esclient.Count(index, esquery.Bool().SetShould(termQueries...))
	if err != nil {
		fmt.Printf("failed to count: %+v\n", err)
		return err
	}
	fmt.Printf("count: %s\n", res.String())

	return nil
}
