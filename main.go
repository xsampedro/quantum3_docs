package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/xsampedro/quantum3_docs/scraper"
)

func main() {
	// Define command line flags
	docType := flag.String("type", "quantum", "Documentation type to scrape: 'quantum' or 'realtime'")
	outputDir := flag.String("output", "", "Output directory (default: output/quantum3 or output/realtime)")
	listOnly := flag.Bool("list-only", false, "Only list URLs without saving content")
	maxConcurrency := flag.Int("concurrency", 5, "Maximum number of concurrent requests")
	
	flag.Parse()

	// Validate docType
	if *docType != "quantum" && *docType != "realtime" {
		fmt.Println("Error: type must be either 'quantum' or 'realtime'")
		flag.Usage()
		os.Exit(1)
	}

	var err error

	// Scrape the documentation based on the specified type
	if *docType == "realtime" {
		fmt.Println("Scraping Photon Realtime documentation...")
		err = scraper.ScrapeRealtimeDocs(*outputDir, *listOnly, *maxConcurrency)
	} else {
		fmt.Println("Scraping Quantum documentation...")
		
		// Create config for Quantum docs
		config := scraper.Config{
			BaseURL:        "https://doc.photonengine.com/quantum/current/getting-started/quantum-intro",
			OutputDir:      *outputDir,
			AllowedDomains: []string{"doc.photonengine.com"},
			ListOnly:       *listOnly,
			BasePath:       "/quantum/current/",
			MaxConcurrency: *maxConcurrency,
			Provider:       scraper.QuantumProvider,
		}

		// Create and run a scraper for Quantum docs
		s := scraper.New(config)
		if err = s.Setup(); err == nil {
			err = s.Run()
		}
	}

	// Check for errors
	if err != nil {
		log.Fatalf("Failed to scrape documentation: %v", err)
	}

	fmt.Println("Documentation scraping completed successfully!")
} 