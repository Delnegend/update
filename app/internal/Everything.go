package check

import (
	"fmt"
	"regexp"
	"strings"
	"update/utils"
)

var everythingLatestVerRe = regexp.MustCompile(`<h2 id="dl">Download Everything ([\d\.]+)</h2>`)

func Everything(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		currVer, err := utils.GetExeVersion(execPath)
		if err != nil {
			return "", fmt.Errorf("can't get executable version: %v", err)
		}
		return strings.TrimSpace(currVer), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestVer, err := func() (string, error) {
		content, err := utils.Fetch("https://www.voidtools.com/downloads")
		if err != nil {
			return "", fmt.Errorf("can't fetch homepage content: %v", err)
		}
		match := everythingLatestVerRe.FindStringSubmatch(string(content))
		if match == nil {
			return "", fmt.Errorf("unexpected website content")
		}
		return match[1], nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	directURL := fmt.Sprintf("https://www.voidtools.com/Everything-%s.x86.zip", latestVer)

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: "https://www.voidtools.com/downloads",
	}
}
