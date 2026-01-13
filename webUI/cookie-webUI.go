package webui

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// AskForCookieWeb attempts to obtain the MoodleSessionag cookie automatically using Chrome,
// falling back to a manual one-click extractor if Chrome is not available.
func AskForCookieWeb() string {
	// Try Chrome automation first
	cookie, err := getCookieWithChrome()
	if err == nil && cookie != "" {
		return cookie
	}

	// Log the error and fall back to manual method
	if err != nil {
		fmt.Printf("Chrome automation not available: %v\n", err)
	}
	fmt.Println("Falling back to manual cookie extraction...")
	openBrowser("https://aulaglobal.uc3m.es")

	return ""
}

// getCookieWithChrome uses chromedp to automate Chrome and capture the cookie after login
func getCookieWithChrome() (string, error) {
	// Check if Chrome/Chromium is available
	if !isChromeAvailable() {
		return "", fmt.Errorf("Chrome or Chromium not found")
	}

	fmt.Println("Please log in with your UC3M credentials. The cookie will be captured automatically.")

	// Create a new Chrome context with visible browser
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.WindowSize(1200, 800),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set a timeout for the entire operation (5 minutes should be enough for login)
	ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	var moodleCookie string

	// Channel to signal when we found the cookie
	cookieFound := make(chan string, 1)

	// Start listening for network events to capture cookies
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if ev, ok := ev.(*network.EventResponseReceived); ok {
			go func() {
				// Get cookies for aulaglobal.uc3m.es
				c, cancel := context.WithTimeout(ctx, 2*time.Second)
				defer cancel()

				var cookies []*network.Cookie
				err := chromedp.Run(c, chromedp.ActionFunc(func(ctx context.Context) error {
					var err error
					cookies, err = network.GetCookies().Do(ctx)
					return err
				}))

				if err != nil {
					return
				}

				for _, cookie := range cookies {
					if cookie.Name == "MoodleSessionag" && strings.Contains(ev.Response.URL, "aulaglobal.uc3m.es") {
						select {
						case cookieFound <- cookie.Value:
						default:
						}
					}
				}
			}()
		}
	})

	// Navigate to AulaGlobal and wait for login
	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate("https://aulaglobal.uc3m.es"),
	)
	if err != nil {
		return "", fmt.Errorf("failed to start Chrome: %w", err)
	}

	// Poll for the authenticated cookie (check if we're on the dashboard/home page after login)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeout := time.After(5 * time.Minute)

	for {
		select {
		case cookie := <-cookieFound:
			// Verify we're actually logged in by checking the current URL
			var currentURL string
			err := chromedp.Run(ctx, chromedp.Location(&currentURL))
			if err == nil {
				// Check if we're past the login page (not on login/index.php and not on SSO)
				if strings.Contains(currentURL, "aulaglobal.uc3m.es") &&
					!strings.Contains(currentURL, "/login/") &&
					!strings.Contains(currentURL, "sso.uc3m.es") && cookie != "" {
					moodleCookie = cookie
					fmt.Println("Cookie captured successfully!")
					return moodleCookie, nil
				}
			}

		case <-ticker.C:
			// Periodically check cookies directly
			var cookies []*network.Cookie
			err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
				var err error
				cookies, err = network.GetCookies().WithURLs([]string{"https://aulaglobal.uc3m.es"}).Do(ctx)
				return err
			}))
			if err != nil {
				continue
			}

			var currentURL string
			if err := chromedp.Run(ctx, chromedp.Location(&currentURL)); err != nil {
				fmt.Printf("Error retrieving current URL: %v\n", err)
				continue
			}

			// Check if we're logged in (on main page, not login page)
			if currentURL == "http://aulaglobal.uc3m.es" &&
				!strings.Contains(currentURL, "/login/") &&
				!strings.Contains(currentURL, "sso.uc3m.es") {
				for _, cookie := range cookies {
					if cookie.Name == "MoodleSessionag" && cookie.Value != "" {
						fmt.Println("Cookie captured successfully!")
						return cookie.Value, nil
					}
				}
			}

		case <-timeout:
			return "", fmt.Errorf("timeout waiting for login")

		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

// isChromeAvailable checks if Chrome or Chromium is available on the system
func isChromeAvailable() bool {
	// chromedp will find Chrome automatically, but let's do a quick check
	browsers := []string{
		"google-chrome",
		"google-chrome-stable",
		"chromium",
		"chromium-browser",
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
	}

	for _, browser := range browsers {
		if _, err := exec.LookPath(browser); err == nil {
			return true
		}
	}

	// On macOS, also check the standard Chrome location
	if runtime.GOOS == "darwin" {
		if out, err := exec.Command("mdfind", "kMDItemCFBundleIdentifier == 'com.google.Chrome'").Output(); err == nil && strings.TrimSpace(string(out)) != "" {
			return true
		}
	}

	return false // Let chromedp handle the detection
}

// openBrowser opens the default browser with the given URL
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		fmt.Println("Error opening browser:", err)
		fmt.Println("Please open the following URL manually:", url)
	}
}
