package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type Crawler struct {
	baseURL        *url.URL
	visitedURLs    sync.Map
	maxDepth       int
	outputDir      string
	concurrency    int
	includeAssets  bool
	skipTLSVerify  bool
	semaphore      chan struct{}
}

func NewCrawler(baseURL string, maxDepth int, outputDir string, concurrency int, includeAssets bool, skipTLSVerify bool) (*Crawler, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	return &Crawler{
		baseURL:        parsedURL,
		maxDepth:       maxDepth,
		outputDir:      outputDir,
		concurrency:    concurrency,
		includeAssets:  includeAssets,
		skipTLSVerify:  skipTLSVerify,
		semaphore:      make(chan struct{}, concurrency),
	}, nil
}

func (c *Crawler) Download(urlStr string, depth int) error {
	if depth > c.maxDepth {
		return nil
	}

	// Проверяем, не скачивали ли мы уже этот URL
	if _, visited := c.visitedURLs.LoadOrStore(urlStr, true); visited {
		return nil
	}

	// Получаем семафор
	c.semaphore <- struct{}{}
	defer func() { <-c.semaphore }()

	// Парсим URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL %s: %v", urlStr, err)
	}

	// Проверяем, что URL относится к тому же домену
	if parsedURL.Host != c.baseURL.Host {
		return nil
	}

	// Создаем HTTP клиент
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: c.skipTLSVerify,
			},
		},
	}

	// Выполняем HTTP запрос
	resp, err := client.Get(urlStr)
	if err != nil {
		return fmt.Errorf("failed to download %s: %v", urlStr, err)
	}
	defer resp.Body.Close()

	// Создаем путь для сохранения файла
	relativePath := strings.TrimPrefix(parsedURL.Path, "/")
	if relativePath == "" {
		relativePath = "index.html"
	}
	filePath := filepath.Join(c.outputDir, relativePath)

	// Создаем необходимые директории
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directories for %s: %v", filePath, err)
	}

	// Создаем файл
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}
	defer file.Close()

	// Копируем содержимое
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save content to %s: %v", filePath, err)
	}

	// Если это HTML документ, ищем ссылки и ресурсы
	if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		// Перечитываем файл для парсинга HTML
		file.Seek(0, 0)
		doc, err := html.Parse(file)
		if err != nil {
			return fmt.Errorf("failed to parse HTML from %s: %v", filePath, err)
		}

		var wg sync.WaitGroup
		c.processNode(doc, depth+1, &wg)
		wg.Wait()
	}

	return nil
}

func (c *Crawler) processNode(n *html.Node, depth int, wg *sync.WaitGroup) {
	if n.Type == html.ElementNode {
		var attr string
		switch n.Data {
		case "a":
			attr = "href"
		case "img", "script", "link":
			if c.includeAssets {
				attr = "src"
				if n.Data == "link" {
					attr = "href"
				}
			}
		}

		if attr != "" {
			for _, a := range n.Attr {
				if a.Key == attr {
					link := a.Val
					if strings.HasPrefix(link, "/") {
						link = c.baseURL.Scheme + "://" + c.baseURL.Host + link
					} else if !strings.HasPrefix(link, "http") {
						base := c.baseURL.String()
						if !strings.HasSuffix(base, "/") {
							base += "/"
						}
						link = base + link
					}

					wg.Add(1)
					go func(url string) {
						defer wg.Done()
						c.Download(url, depth)
					}(link)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		c.processNode(c, depth, wg)
	}
}

func main() {
	url := flag.String("url", "", "URL to download")
	depth := flag.Int("depth", 1, "Maximum depth for recursive downloading")
	outputDir := flag.String("output", "downloaded", "Output directory")
	concurrency := flag.Int("concurrency", 5, "Number of concurrent downloads")
	includeAssets := flag.Bool("assets", true, "Download assets (images, scripts, styles)")
	skipTLSVerify := flag.Bool("skip-tls-verify", false, "Skip TLS certificate verification")
	
	flag.Parse()

	if *url == "" {
		fmt.Println("Please provide a URL using -url flag")
		flag.PrintDefaults()
		os.Exit(1)
	}

	crawler, err := NewCrawler(*url, *depth, *outputDir, *concurrency, *includeAssets, *skipTLSVerify)
	if err != nil {
		fmt.Printf("Error creating crawler: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting download of %s (depth: %d)\n", *url, *depth)
	if err := crawler.Download(*url, 0); err != nil {
		fmt.Printf("Error during download: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Download completed successfully")
}