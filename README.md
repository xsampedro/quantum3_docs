# Quantum3 API Documentation Scraper

A simple Go application that scrapes API documentation from a vendor's website and saves it locally.

## Purpose

This repository was created to facilitate uploading API documentation to [Context7](https://context7.com/add-library), which requires documentation to be in markdown format and hosted in a GitHub repository.

## Prerequisites

- Go 1.21 or later
- Internet connection to access the vendor's documentation

## Project Structure

```
.
├── README.md           # Project documentation
├── go.mod              # Go module definition
├── main.go             # Main application code
└── output/             # Root output directory
    └── quantum3/       # Output directory for Quantum3 documentation
```

## Getting Started

1. Clone the repository:
   ```
   git clone https://github.com/xsampedro/quantum3_docs.git
   cd quantum3_docs
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Run the application:
   ```
   go run main.go -list=false
   ```

## GitHub Actions

This repository includes a GitHub Action workflow that can be manually triggered to re-run the scraper:

1. Go to the "Actions" tab in your GitHub repository
2. Select the "Manual API Docs Scraper" workflow
3. Click "Run workflow"
4. Optionally provide a reason for the update
5. Click "Run workflow" again to start the process

The action will run the scraper, commit any changes to the documentation, and push them to the repository.

## Configuration

To configure the application for your specific needs:

1. Update the vendor domain in `main.go` by modifying the `domain` flag
2. Set the correct API documentation starting URL by updating the `baseURL` flag
3. Change the output directory using the `output` flag (defaults to `output/quantum3`)
4. Use `-list=true` to only list URLs without downloading (defaults to `true`)
5. Adjust concurrency with the `concurrency` flag (defaults to twice the number of CPU cores)

Example with custom settings:
```
go run main.go -url="https://example.com/docs" -output="output/custom" -list=false -concurrency=4
```

## License

MIT License

Copyright (c) 2023 xsampedro

This license applies only to the scraper code itself. The content in the 'output' folder consists of API documentation scraped from the vendor's website and remains subject to the vendor's own copyright and licensing terms. No copyright claim is made over this scraped content.

For full license text, see the [LICENSE](LICENSE) file. 