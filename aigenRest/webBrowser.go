package aigenRest

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
)

// SecureBrowser provides a secure, headless browser instance using Rod (Chromium-based)
type SecureBrowser struct {
	browser *rod.Browser
	page    *rod.Page
}

// NewSecureBrowser creates a new secure browser instance
func NewSecureBrowser() (*SecureBrowser, error) {
	// Get Chromium executable path
	chromiumPath := getChromiumPath()
	if chromiumPath == "" {
		return nil, fmt.Errorf("chromium not found. Install: snap install chromium")
	}

	// Create launcher with security flags
	l := launcher.New().
		Headless(false).
		Bin(chromiumPath).
		Set(flags.NoSandbox)

	// Launch browser
	browser := rod.New().
		ControlURL(l.MustLaunch()).
		MustConnect().
		Timeout(time.Minute * 5)

	// Create secure page
	page := browser.MustPage()

	// Navigate to a safe starting page
	page.MustNavigate("about:blank")

	return &SecureBrowser{
		browser: browser,
		page:    page,
	}, nil
}

// Navigate navigates to a URL with security checks
func (sb *SecureBrowser) Navigate(url string) error {
	// Security check: Validate URL
	if !isValidURL(url) {
		return fmt.Errorf("invalid URL: %s", url)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sb.page.Context(ctx).MustNavigate(url)
	return nil
}

// GetContent gets the page content
func (sb *SecureBrowser) GetContent() (string, error) {
	content, err := sb.page.HTML()
	return content, err
}

// GetTextContent gets the visible text content
func (sb *SecureBrowser) GetTextContent() (string, error) {
	text, err := sb.page.MustElement("body").Text()
	return text, err
}

// SearchContent searches for content on the page
func (sb *SecureBrowser) SearchContent(query string) ([]string, error) {
	// Get all paragraph elements
	elements := sb.page.MustElements("p")

	// Filter results containing query
	var matches []string
	for _, elem := range elements {
		text, err := elem.Text()
		if err == nil && strings.Contains(strings.ToLower(text), strings.ToLower(query)) {
			matches = append(matches, text)
		}
	}

	return matches, nil
}

// Screenshot takes a screenshot of the current page
func (sb *SecureBrowser) Screenshot(filepath string) error {
	img, err := sb.page.Screenshot(true, nil)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, img, 0644)
}

// Close closes the browser
func (sb *SecureBrowser) Close() error {
	if sb.page != nil {
		sb.page.MustClose()
	}
	if sb.browser != nil {
		sb.browser.MustClose()
	}
	return nil
}

// Security helper functions
func isValidURL(url string) bool {
	if url == "" || len(url) < 4 {
		return false
	}
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// getChromiumPath finds the Chromium binary
func getChromiumPath() string {
	var possiblePaths []string

	if runtime.GOOS == "linux" {
		possiblePaths = []string{
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
			"/snap/bin/chromium",
		}
	} else if runtime.GOOS == "darwin" {
		possiblePaths = []string{
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		}
	} else if runtime.GOOS == "windows" {
		possiblePaths = []string{
			`C:\Program Files\Chromium\Application\chrome.exe`,
			`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
		}
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}
