# Quantum3 and Realtime API Documentation Scraper

A simple Go application that scrapes API documentation from Photon's websites and saves it locally in Markdown format.

## Purpose

This repository was created to facilitate uploading API documentation to [Context7](https://context7.com/add-library), which requires documentation to be in markdown format and hosted in a GitHub repository.

## Features

- Scrapes Quantum 3 and Photon Realtime documentation from the official website
- Converts HTML to clean Markdown format
- Properly handles and formats code blocks
- Preserves original document structure
- Organizes documents in a logical folder structure

## Prerequisites

- Go 1.21 or later
- Internet connection to access the vendor's documentation

## Project Structure

```
.
├── README.md           # Project documentation
├── go.mod              # Go module definition
├── main.go             # Main application code
├── scraper/            # Scraper package code
│   └── scraper.go      # Core scraping functionality
└── output/             # Root output directory
    ├── quantum3/       # Output directory for Quantum3 documentation
    └── realtime/       # Output directory for Photon Realtime documentation
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

3. Run the application to scrape Quantum documentation:
   ```
   go run main.go -type=quantum -list-only=false
   ```

4. Run the application to scrape Realtime documentation:
   ```
   go run main.go -type=realtime -list-only=false
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

1. Use the `-type` flag to specify which documentation to scrape:
   - `quantum`: Scrapes Quantum documentation
   - `realtime`: Scrapes Photon Realtime documentation

2. Change the output directory using the `-output` flag (defaults to `output/quantum3` for Quantum or `output/realtime` for Realtime)

3. Use `-list-only=true` to only list URLs without downloading (defaults to `false`)

4. Adjust concurrency with the `-concurrency` flag (defaults to 5)

Example with custom settings:
```
go run main.go -type=realtime -output="output/custom-realtime" -list-only=false -concurrency=4
```

## License

MIT License

Copyright (c) 2025 xsampedro.com

This license applies only to the scraper code itself. The content in the 'output' folder consists of API documentation scraped from the vendor's website and remains subject to the vendor's own copyright and licensing terms. No copyright claim is made over this scraped content.

For full license text, see the [LICENSE](LICENSE) file. 
