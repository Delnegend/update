package check

import (
	"fmt"
	"regexp"
	"strings"
	"update/utils"
)

var handbrakeLatestVerRe = regexp.MustCompile(`HandBrake ([\d\.]+)`)

func Handbrake(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		currVer, err := utils.GetExeVersion(execPath)
		if err != nil {
			return "", fmt.Errorf("can't get executable version: %v", err)
		}
		return strings.TrimSuffix(currVer, ".0"), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestVer, err := func() (string, error) {
		content, err := utils.Fetch("https://handbrake.fr/")
		if err != nil {
			return "", fmt.Errorf("can't fetch the homepage: %v", err)
		}
		match := handbrakeLatestVerRe.FindStringSubmatch(string(content))
		if match == nil {
			return "", fmt.Errorf("unexpected website content")
		}
		return strings.TrimSpace(match[1]), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	directURL := fmt.Sprintf("https://github.com/HandBrake/HandBrake/releases/download/%s/HandBrake-%s-x86_64-Win_GUI.zip", latestVer, latestVer)

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: "https://handbrake.fr/downloads.php",
	}
}
