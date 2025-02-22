package cookies

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

const URLToken = "https://aulaglobal.uc3m.es/admin/tool/mobile/launch.php?service=moodle_mobile_app&passport=82.93261629596182&urlscheme=moodlemobile"

func CookieToToken(cookie string) (string, error) {
	// This function will convert the cookie to the token
	_, err := getToken(cookie)
	if err == nil {
		return "", err
	}
	token, shouldNotReturn := extractTokenFromError(err)
	if shouldNotReturn {
		return "", errors.New("failed to extract token from error")
	}
	return token, nil
}

func getToken(cookie string) (string, error) {
	// Set the URL and headers
	client := &http.Client{}
	req, err := http.NewRequest("GET", URLToken, nil)
	if err != nil {
		return "", err
	}

	// Add Cookie header
	cookieValue := fmt.Sprintf("MoodleSessionag=%s", cookie)
	req.Header.Add("Cookie", cookieValue)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Match token in the response
	re := regexp.MustCompile(`token=([^&]+)$`)
	matches := re.FindStringSubmatch(string(body))
	if matches == nil || len(matches) < 2 {
		return "", errors.New("failed to find token in response")
	}

	token := matches[1]
	// Remove everything from `=` to the end
	parts := strings.Split(token, ")")
	if len(parts) < 1 {
		return "", errors.New("failed to parse token")
	}

	// Decode base64 token
	decodedToken, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return "", errors.New("failed to decode base64 token")
	}

	decodedTokenStrComplete := string(decodedToken)
	decodedParts := strings.Split(decodedTokenStrComplete, ":::")
	if len(decodedParts) < 2 {
		return "", errors.New("failed to parse decoded token")
	}

	decodedTokenStr := decodedParts[1]
	return decodedTokenStr, nil
}

func extractTokenFromError(err error) (string, bool) {
	// Aplpy a regex to get the token: r"token=([^&]+)$"
	// Decode the token from base64
	// Split the token by ":::" and get the second part
	// Print the token
	// Convert from byte array to string and split by ":::"
	pattern := regexp.MustCompile(`token=([^&]+)$`)
	matches := pattern.FindStringSubmatch(err.Error())

	if matches == nil {
		return "", true
	}

	token := matches[1]
	token = strings.Replace(token, ")", "", -1)
	parts := strings.Split(token, "\"")
	if len(parts) < 1 {
		fmt.Println("No token found")
		return "", true
	}

	decodedToken, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		fmt.Println("Token:", err)
		return "", true
	}

	decodedTokenList := strings.Split(string(decodedToken), ":::")
	if len(decodedTokenList) != 2 {
		fmt.Println("\":::\" not found in the token")
		return "", true
	}

	token = decodedTokenList[1]
	return token, false
}
