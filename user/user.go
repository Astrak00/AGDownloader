package user

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	types "github.com/Astrak00/AGDownloader/types"

	"github.com/fatih/color"
)

/*
Gets the userID necessary to get the courses
TODO: Change this to a json response
*/
func GetUserInfo(token string) (types.UserInfo, error) {
	url := fmt.Sprintf("https://%s%s?wstoken=%s&wsfunction=core_webservice_get_site_info", types.Domain, types.Webservice, token)
	resp, err := http.Get(url)
	if err != nil {
		return types.UserInfo{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.UserInfo{}, err
	}

	if strings.Contains(string(body), "invalidtoken") {
		return types.UserInfo{}, fmt.Errorf("invalid token")
	}

	var userInfo types.UserInfo

	// Find the fullname key and value
	fullName := regexp.MustCompile(`<KEY name="fullname"><VALUE>([^<]+)</VALUE>`)
	matches := fullName.FindStringSubmatch(string(body))
	if len(matches) >= 1 {
		userInfo.FullName = matches[1]
	} else {
		color.Red("Fullname not found\n")
	}

	// Find the userid key and value
	userID := regexp.MustCompile(`<KEY name="userid"><VALUE>([^<]+)</VALUE>`)
	matches = userID.FindStringSubmatch(string(body))
	if len(matches) > 1 {
		userInfo.UserID = matches[1]
	} else {
		color.Red("UserID not found\n")
	}

	return userInfo, nil
}
