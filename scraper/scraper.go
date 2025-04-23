package scraper

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

// Config contains the configuration for the API docs scraper
type Config struct {
	BaseURL        string
	OutputDir      string
	AllowedDomains []string
	ListOnly       bool
	BasePath       string
	MaxConcurrency int
}

// Scraper manages the scraping of API documentation
type Scraper struct {
	config    Config
	collector *colly.Collector
	visited   map[string]bool
	mutex     sync.Mutex
	converter *md.Converter
}

// New creates a new API documentation scraper
func New(config Config) *Scraper {
	if config.OutputDir == "" {
		config.OutputDir = "output/quantum3"
	}

	if config.BasePath == "" {
		config.BasePath = "/quantum/current/"
	}

	if config.MaxConcurrency == 0 {
		config.MaxConcurrency = 10 // Default concurrency
	}

	// Create a new markdown converter with default options
	converter := md.NewConverter("", true, nil)

	// Add GitHub flavored markdown plugins
	converter.Use(plugin.GitHubFlavored())
	
	// Add custom rules to handle specific HTML elements better
	converter.AddRules(
		md.Rule{
			Filter: []string{"h1", "h2", "h3", "h4", "h5", "h6"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				level := selec.Get(0).Data[1] // Get the heading level from the tag name
				result := strings.Repeat("#", int(level-'0')) + " " + strings.TrimSpace(content) + "\n\n"
				return &result
			},
		},
		md.Rule{
			Filter: []string{"code", "pre"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				// Determine the language if possible
				lang := selec.AttrOr("class", "")
				if strings.Contains(lang, "language-") {
					lang = strings.Replace(lang, "language-", "", 1)
				} else {
					lang = ""
				}
				
				// Check if content already contains triple backticks
				if strings.Contains(content, "```") {
					// If the content already has backticks, we need to ensure it's not wrapped again
					return &content
				}
				
				result := "```" + lang + "\n" + content + "\n```\n\n"
				return &result
			},
		},
	)

	return &Scraper{
		config: config,
		collector: colly.NewCollector(
			colly.AllowedDomains(config.AllowedDomains...),
			colly.Async(true),
			colly.MaxDepth(5),
		),
		visited:   make(map[string]bool),
		converter: converter,
	}
}

// shouldProcessURL determines if the URL should be processed
func (s *Scraper) shouldProcessURL(url string) bool {
	// Skip non-English docs
	if strings.Contains(url, "/ja-jp/") || strings.Contains(url, "/ko-kr/") || 
	   strings.Contains(url, "/zh-cn/") || strings.Contains(url, "/zh-tw/") {
		return false
	}
	
	// Only process URLs in the specified base path
	return strings.Contains(url, s.config.BasePath)
}

// generateOutputPath creates a valid output file path for a URL
func (s *Scraper) generateOutputPath(url string) string {
	// Extract path from URL
	urlPath := strings.TrimPrefix(url, "https://"+s.config.AllowedDomains[0])
	
	// Remove the base path prefix for cleaner filenames
	relativePath := strings.TrimPrefix(urlPath, s.config.BasePath)
	if relativePath == "" || relativePath == "/" {
		relativePath = "index"
	}

	// Clean up path and ensure it ends with .md
	if strings.HasSuffix(relativePath, "/") {
		relativePath = strings.TrimSuffix(relativePath, "/")
		if filepath.Ext(relativePath) == "" {
			relativePath += "/index"
		}
	}
	
	// Create full path in output directory
	dirPath := filepath.Join(s.config.OutputDir, filepath.Dir(relativePath))
	filename := filepath.Base(relativePath)
	
	// Ensure filename has .md extension
	if filepath.Ext(filename) != "" {
		filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	}
	
	return filepath.Join(dirPath, filename+".md")
}

// Setup configures the scraper
func (s *Scraper) Setup() error {
	// Set concurrency limit
	s.collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: s.config.MaxConcurrency,
	})

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(s.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Configure the collector
	s.collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(link)

		if s.shouldProcessURL(absoluteURL) {
			s.mutex.Lock()
			alreadyVisited := s.visited[absoluteURL]
			if !alreadyVisited {
				s.visited[absoluteURL] = true
				s.mutex.Unlock()
				if s.config.ListOnly {
					fmt.Println("Found:", absoluteURL)
				}
				e.Request.Visit(link)
			} else {
				s.mutex.Unlock()
			}
		}
	})

	if !s.config.ListOnly {
		s.collector.OnResponse(func(r *colly.Response) {
			// Generate output path
			outputPath := s.generateOutputPath(r.Request.URL.String())
			
			// Create directory if needed
			dirPath := filepath.Dir(outputPath)
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				log.Printf("Error creating directory %s: %v", dirPath, err)
				return
			}
			
			// Convert HTML to Markdown
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(r.Body)))
			if err != nil {
				log.Printf("Error parsing HTML from %s: %v", r.Request.URL, err)
				return
			}
			
			// Extract main content and remove navigation, header, footer
			// This is site-specific and may need adjustment
			var mainContent string
			doc.Find("body").Each(func(i int, sel *goquery.Selection) {
				// Remove navigation, header, footer elements
				sel.Find("nav, header, footer, script, style, .sidebar, #sidebar").Remove()
				
				// Find main content - typical container for article content
				content := sel.Find("main, article, .content, .main-content, #content")
				if content.Length() > 0 {
					// Use the content we found
					htmlContent, err := content.Html()
					if err == nil {
						mainContent = htmlContent
					}
				} else {
					// Fallback to body if no specific content container found
					htmlContent, err := sel.Html()
					if err == nil {
						mainContent = htmlContent
					}
				}
			})
			
			// Convert to markdown
			markdown, err := s.converter.ConvertString(mainContent)
			if err != nil {
				log.Printf("Error converting HTML to Markdown for %s: %v", r.Request.URL, err)
				
				// Fallback: save raw HTML
				if err := os.WriteFile(outputPath, r.Body, 0644); err != nil {
					log.Printf("Error saving raw HTML file %s: %v", outputPath, err)
				} else {
					fmt.Printf("Saved raw HTML (conversion failed): %s\n", outputPath)
				}
				return
			}
			
			// Add URL reference at the top of the file
			markdown = fmt.Sprintf("# %s\n\n_Source: %s_\n\n%s", 
				filepath.Base(strings.TrimSuffix(outputPath, ".md")),
				r.Request.URL.String(),
				markdown)
			
			// Save the Markdown content
			if err := os.WriteFile(outputPath, []byte(markdown), 0644); err != nil {
				log.Printf("Error saving markdown file %s: %v", outputPath, err)
			} else {
				fmt.Printf("Saved: %s\n", outputPath)
			}
		})
	}

	s.collector.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s failed with error: %v", r.Request.URL, err)
	})

	return nil
}

// Run starts the scraping process
func (s *Scraper) Run() error {
	s.mutex.Lock()
	s.visited[s.config.BaseURL] = true
	s.mutex.Unlock()
	
	if s.config.ListOnly {
		fmt.Println("Starting URL:", s.config.BaseURL)
	}
	
	err := s.collector.Visit(s.config.BaseURL)
	if err != nil {
		return err
	}
	
	// Wait for all concurrent requests to finish
	s.collector.Wait()
	
	// If we're saving files, let's clean up any markdown issues
	if !s.config.ListOnly {
		s.cleanupMarkdownFiles()
	}
	
	return nil
}

// cleanupMarkdownFiles fixes common issues in generated markdown files
func (s *Scraper) cleanupMarkdownFiles() error {
	// Walk through the output directory
	return filepath.Walk(s.config.OutputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Only process markdown files
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			// Read file
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			
			// Fix nested code blocks - replace triple backticks inside code blocks
			markdown := string(content)
			
			// Replace pattern: ```(language)\n```(language2)\n with just the inner block
			markdown = cleanNestedCodeBlocks(markdown)
			
			// Save the cleaned content
			if err := os.WriteFile(path, []byte(markdown), 0644); err != nil {
				return err
			}
			
			fmt.Printf("Cleaned: %s\n", path)
		}
		
		return nil
	})
}

// cleanNestedCodeBlocks fixes nested code blocks in markdown
func cleanNestedCodeBlocks(content string) string {
	// Replace pattern like: "```language\n```language2\n" with "```language2\n"
	nestedCodeBlockRegex := regexp.MustCompile("```(.+?)\n```(.+?)\n")
	for nestedCodeBlockRegex.MatchString(content) {
		content = nestedCodeBlockRegex.ReplaceAllString(content, "```$2\n")
	}
	
	// Replace pattern like: "```\n```language\n" with "```language\n"
	nestedEmptyCodeBlockRegex := regexp.MustCompile("```\n```(.+?)\n")
	for nestedEmptyCodeBlockRegex.MatchString(content) {
		content = nestedEmptyCodeBlockRegex.ReplaceAllString(content, "```$1\n")
	}
	
	// Replace any `````` double backticks with single ```
	content = strings.ReplaceAll(content, "``````", "```")
	
	return content
}

// GetVisitedURLs returns a list of all visited URLs
func (s *Scraper) GetVisitedURLs() []string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	urls := make([]string, 0, len(s.visited))
	for url := range s.visited {
		urls = append(urls, url)
	}
	return urls
} 