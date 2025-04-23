# Quantum3 API Documentation Scraper

A simple Go application that scrapes API documentation from a vendor's website and saves it locally.

## Prerequisites

- Go 1.21 or later
- Internet connection to access the vendor's documentation

## Project Structure

```
.
├── README.md           # Project documentation
├── go.mod              # Go module definition
├── main.go             # Main application code
└── docs/               # Output directory for scraped documentation
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
   go run main.go
   ```

## Configuration

To configure the application for your specific needs:

1. Update the vendor domain in `main.go` by modifying the `AllowedDomains` parameter
2. Set the correct API documentation starting URL by updating the `apiDocsURL` variable
3. Customize the scraping logic in the `OnHTML` handler to match the structure of the vendor's documentation

## License

[Add license information here] 