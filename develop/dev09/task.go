package main

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	// "path"
	"path/filepath"
	"strings"
)

// downloadFile downloads a file from the given URL and saves it to the specified path.
func downloadFile(url string, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// parseHTML parses the HTML content and extracts all resource URLs.
func parseHTML(baseURL *url.URL, body io.Reader) ([]string, error) {
	var resourceURLs []string
	z := html.NewTokenizer(body)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return resourceURLs, nil
			}
			return nil, z.Err()
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			for _, attr := range t.Attr {
				if attr.Key == "src" || attr.Key == "href" {
					resURL, err := url.Parse(attr.Val)
					if err != nil {
						continue
					}
					if !resURL.IsAbs() {
						resURL = baseURL.ResolveReference(resURL)
					}
					resourceURLs = append(resourceURLs, resURL.String())
				}
			}
		}
	}
}

// createDirIfNotExist creates a directory if it does not already exist.
func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func sanitizeFileName(fileName string) string {
	return strings.ReplaceAll(fileName, ":", "_")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No args provided")
		fmt.Println("Usage: ./task URL")
		os.Exit(0)
	}

	source := os.Args[1]
	parsedURL, err := url.Parse(source)
	if err != nil {
		log.Fatalf("%s is not a valid link\n", source)
	}
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "http"
		source = parsedURL.String()
	}

	resp, err := http.Get(source)
	if err != nil {
		log.Fatalf("Error downloading the main page: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to download the page: %s\n", resp.Status)
	}

	// Create a base directory for the website
	baseDir := sanitizeFileName(parsedURL.Host)
	if err := createDirIfNotExist(baseDir); err != nil {
		log.Fatalf("Failed to create base directory: %v\n", err)
	}

	// Save the main HTML file
	mainHTMLFile := filepath.Join(baseDir, "index.html")
	out, err := os.Create(mainHTMLFile)
	if err != nil {
		log.Fatalf("Failed to create main HTML file: %v\n", err)
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("Failed to save main HTML file: %v\n", err)
	}

	// Re-download the HTML for parsing
	resp, err = http.Get(source)
	if err != nil {
		log.Fatalf("Error re-downloading the main page for parsing: %v\n", err)
	}
	defer resp.Body.Close()

	// Parse the HTML and download resources
	resourceURLs, err := parseHTML(parsedURL, resp.Body)
	if err != nil {
		log.Fatalf("Error parsing HTML: %v\n", err)
	}

	for _, resURL := range resourceURLs {
		u, err := url.Parse(resURL)
		if err != nil {
			log.Printf("Skipping invalid resource URL %s: %v\n", resURL, err)
			continue
		}

		// Create directory structure for the resource
		resourcePath := filepath.Join(baseDir, u.Host, u.Path)
		resourceDir := filepath.Dir(resourcePath)
		if err := createDirIfNotExist(resourceDir); err != nil {
			log.Printf("Failed to create directory for resource %s: %v\n", resURL, err)
			continue
		}

		// Download and save the resource
		log.Printf("Downloading resource %s\n", resURL)
		if err := downloadFile(resURL, resourcePath); err != nil {
			log.Printf("Failed to download resource %s: %v\n", resURL, err)
		}
	}

	log.Println("Website downloaded successfully.")
}
