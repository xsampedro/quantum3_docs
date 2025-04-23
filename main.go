package main

import (
	"flag"
	"log"
	"runtime"
	"sort"

	"github.com/xsampedro/quantum3_docs/scraper"
)

func main() {
	// Define command-line flags
	baseURL := flag.String("url", "https://doc.photonengine.com/quantum/current/quantum-intro", "Base URL for the Quantum documentation")
	outputDir := flag.String("output", "output", "Directory to save the scraped documentation")
	domain := flag.String("domain", "doc.photonengine.com", "Allowed domain for scraping")
	listOnly := flag.Bool("list", true, "Only list URLs without downloading content")
	basePath := flag.String("base-path", "/quantum/current/", "Base path to restrict scraping")
	concurrency := flag.Int("concurrency", runtime.NumCPU()*2, "Maximum number of concurrent requests")
	flag.Parse()

	// Create and configure the scraper
	config := scraper.Config{
		BaseURL:        *baseURL,
		OutputDir:      *outputDir,
		AllowedDomains: []string{*domain},
		ListOnly:       *listOnly,
		BasePath:       *basePath,
		MaxConcurrency: *concurrency,
	}

	apiScraper := scraper.New(config)
	if err := apiScraper.Setup(); err != nil {
		log.Fatalf("Failed to set up scraper: %v", err)
	}

	// Run the scraper
	log.Printf("Starting to scan Quantum documentation from %s", *baseURL)
	log.Printf("Listing mode: %v, Concurrency: %d", *listOnly, *concurrency)
	if !*listOnly {
		log.Printf("Output directory: %s", *outputDir)
		log.Printf("Files will be saved in Markdown format")
	}
	
	if err := apiScraper.Run(); err != nil {
		log.Fatalf("Error during scraping: %v", err)
	}

	// If in list mode, print a summary of all found URLs
	if *listOnly {
		urls := apiScraper.GetVisitedURLs()
		sort.Strings(urls)
		log.Printf("\nFound %d pages in total.\n", len(urls))
		log.Println("To download these pages, run the tool with -list=false")
	}
} 