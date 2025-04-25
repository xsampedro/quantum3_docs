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

	// Determine configuration based on docType
	var config scraper.Config

	if *docType == "realtime" {
		fmt.Println("Scraping Photon Realtime documentation...")
		config = scraper.Config{
			BaseURL:           "https://doc.photonengine.com/realtime/current/getting-started/realtime-intro",
			OutputDir:         *outputDir,
			AllowedDomains:    []string{"doc.photonengine.com"},
			ListOnly:          *listOnly,
			BasePath:          "/realtime/current/",
			MaxConcurrency:    *maxConcurrency,
			ContentSelectors:  "main, article, .content, .main-content, #content, div[role='main'], body",
			ExcludeURLPattern: []string{"/ja-jp/", "/ko-kr/", "/zh-cn/", "/zh-tw/"},
		}

		// Set default output directory if not specified
		if config.OutputDir == "" {
			config.OutputDir = "output/realtime"
		}
	} else {
		fmt.Println("Scraping Quantum documentation...")
		config = scraper.Config{
			BaseURL:           "https://doc.photonengine.com/quantum/current/getting-started/quantum-intro",
			OutputDir:         *outputDir,
			AllowedDomains:    []string{"doc.photonengine.com"},
			ListOnly:          *listOnly,
			BasePath:          "/quantum/current/",
			MaxConcurrency:    *maxConcurrency,
			ContentSelectors:  "main, article, .content, .main-content, #content",
			ExcludeURLPattern: []string{"/ja-jp/", "/ko-kr/", "/zh-cn/", "/zh-tw/"},
		}

		// Set default output directory if not specified
		if config.OutputDir == "" {
			config.OutputDir = "output/quantum3"
		}
	}

	// Create and run the scraper
	s := scraper.New(config)
	if err = s.Setup(); err == nil {
		err = s.Run()
	}

	// Check for errors
	if err != nil {
		log.Fatalf("Failed to scrape documentation: %v", err)
	}

	fmt.Println("Documentation scraping completed successfully!")
}
