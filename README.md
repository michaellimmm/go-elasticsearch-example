# Go Elasticsearch Example

==========================

## Overview

This is a simple example of using Elasticsearch with Go. The project provides a command-line interface to interact with an Elasticsearch instance.

## Prerequisites

- Go (version 1.17 or later)
- Elasticsearch (version 7.x or later)

## How to Use

### Building the Project

To build the project, run the following command:

```bash
go build cmd/cli/main.go
```

This will create an executable file named `main` in the `cmd/cli` directory.

### Running the Project

To run the project, use the following command:

```bash
go run cmd/cli/main.go -command <command> -file <file>
```

Replace `<command>` with one of the following options:

- `create-index`: creates an index in Elasticsearch
- `indexing`: indexes a CSV file in Elasticsearch
- `match-docs`: searches for documents in Elasticsearch

Replace `<file>` with the path to a CSV file (required for `indexing` and `match-docs` commands).

### Example Usage

To match documents in Elasticsearch using a CSV file:

```bash
go run cmd/cli/main.go -command match-docs -file ./bucket/sample2.csv
```

### Configuration

The project uses the following configuration:

- Elasticsearch instance: `http://localhost:9200`
- Index name: `item_index_ja` (or `item_index_en` for English language)

You can modify these settings by editing the `internal/config/index/index.go` file.

## Contributing

Contributions are welcome! Please submit a pull request with your changes.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
