package check

import (
	"fmt"
	"regexp"
	"strings"
	"update/utils"
)

var topazGigapixelAILatestVerRe = regexp.MustCompile(`Gigapixel v([\d\.]+)`)

func TopazGigapixelAI(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		currVer, err := utils.GetExeVersion(execPath)
		if err != nil {
			return "", err
		}
		return strings.Join(strings.Split(currVer, ".")[:3], "."), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestVer, err := func() (string, error) {
		content, err := utils.Fetch("https://community.topazlabs.com/c/gigapixel-ai/gigapixel-ai/")
		if err != nil {
			return "", fmt.Errorf("can't fetch the homepage: %v", err)
		}
		match := topazGigapixelAILatestVerRe.FindStringSubmatch(string(content))
		if match == nil {
			return "", fmt.Errorf("unexpected website content")
		}
		return strings.TrimSpace(match[1]), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DownloadPageURL: "https://www.topazlabs.com/gigapixel",
	}
}
