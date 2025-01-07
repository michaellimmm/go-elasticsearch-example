# **Go Elasticsearch Example**

## Overview

This project provides a simple example of using Elasticsearch with Go. It includes a command-line interface to interact with an Elasticsearch instance.

## Prerequisites

- Go (version 1.17 or later)
- Docker
- Docker Compose

## Running the Project

### Step 1: Start Elasticsearch using Docker Compose

Before running the command-line interface, you need to start the Elasticsearch instance using Docker Compose. Run the following command:

```bash
docker-compose up -d
```

This will start the Elasticsearch instance in detached mode.

### Step 2: Build and Run the Command-Line Interface

Once the Elasticsearch instance is running, you can build and run the command-line interface using the following command:

```bash
go build cmd/cli/main.go
```

This will create an executable file named `main` in the `cmd/cli` directory.

### Step 3: Run the Command-Line Interface

To run the command-line interface, use the following command:

```bash
go run cmd/cli/main.go -command <command> -file <file>
```

Replace `<command>` with one of the following options:

- `create-index`: creates an index in Elasticsearch
- `indexing`: indexes a CSV file in Elasticsearch
- `match-docs`: searches for documents in Elasticsearch
- `upload-file-to-gcs`: uploads a file to Google Cloud Storage

## Generating CSV Files

The `cmd/generatecsv` command generates a CSV file containing sample data. You can run it using the following command:

```bash
go run cmd/generatecsv/main.go
```

This will generate a CSV file named `sample.csv` in the current directory.

## Converting CSV to Parquet

The `cmd/csvtoparquet` command converts a CSV file to a Parquet file. You can run it using the following command:

```bash
go run cmd/csvtoparquet/main.go -input <input_csv_file> -output <output_parquet_file>
```

Replace `<input_csv_file>` with the path to the input CSV file and `<output_parquet_file>` with the desired path for the output Parquet file.

## Configuration

The project uses the following configuration:

- Elasticsearch instance: `http://localhost:9200`
- Index name: `item_index_ja` (or `item_index_en` for English language)

You can modify these settings by editing the `internal/config/index/index.go` file.

## Contributing

Contributions are welcome! Please submit a pull request with your changes.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

Let me know if this looks good or if you'd like me to make any changes!
