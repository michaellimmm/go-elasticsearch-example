package main

import (
	"flag"
	"fmt"
	"github/shaolim/go-elasticsearch-example/internal/config"
	"github/shaolim/go-elasticsearch-example/internal/usecase"
	"github/shaolim/go-elasticsearch-example/pkg/esclient"
)

type Command string

const (
	CreateIndex Command = "create-index"
	Indexing    Command = "indexing"
	MatchDocs   Command = "match-docs"
)

func main() {
	var (
		command      string
		filename     string
		languageCode string
	)
	help := flag.Bool("help", false, "Display help")

	flag.StringVar(&command, "command", "", "Command eg. create-index, indexing, match-docs")
	flag.StringVar(&filename, "file", "", "File name")
	flag.StringVar(&languageCode, "lang", "ja", "Language code")

	flag.Parse()

	if *help || command == "" {
		flag.PrintDefaults()
		return
	}

	switch Command(command) {
	case CreateIndex:
		if err := createIndex(); err != nil {
			fmt.Println(err)
		}
	case Indexing:
		if filename == "" {
			fmt.Println("filename is required to run this indexing command")
			return
		}
		if err := indexing(languageCode, filename); err != nil {
			fmt.Println(err)
		}
	case MatchDocs:
		if filename == "" {
			fmt.Println("filename is required to run this match-docs command")
			return
		}

		if err := matchDocs(filename, languageCode); err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Printf("unknown command: %s, valid commands: create-index, indexing, match-docs\n", command)
	}
}

func createIndex() error {
	client := esclient.NewClient("http://localhost:9200")

	createIndexUC := usecase.NewCreateIndexUseCase(client)
	if err := createIndexUC.Execute(); err != nil {
		fmt.Printf("failed to create index, error: %v\n", err)
		return err
	}

	return nil
}

func indexing(languageCode string, filename string) error {
	client := esclient.NewClient("http://localhost:9200")

	index := config.ItemIndexJa
	if languageCode == "en" {
		index = config.ItemIndexEn
	}

	indexingUC := usecase.NewDocsInsertUseCase(client)
	if err := indexingUC.Execute(index, filename); err != nil {
		fmt.Printf("failed to indexing, error: %v\n", err)
		return err
	}
	return nil
}

func matchDocs(filename string, languageCode string) error {
	client := esclient.NewClient("http://localhost:9200")
	index := config.ItemIndexJa
	if languageCode == "en" {
		index = config.ItemIndexEn
	}

	matchDocs := usecase.NewSampleDocsUseCase(client)
	if err := matchDocs.Execute(index, filename); err != nil {
		fmt.Printf("failed to match docs, error: %v\n", err)
		return err
	}
	return nil
}
