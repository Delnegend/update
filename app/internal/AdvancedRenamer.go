package check

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"update/utils"
)

var (
	arenLatestVerRe = regexp.MustCompile(`Version: (\d+\.\d+)`)
	arenCurrVerRe   = regexp.MustCompile(`Advanced Renamer (\d+\.\d+)`)
)

func AdvancedRenamer(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		currVerBytes, err := exec.Command(execPath, `-h`).Output()
		if err != nil {
			return "", err
		}
		match := arenCurrVerRe.FindStringSubmatch(string(currVerBytes))
		if match == nil {
			return "", fmt.Errorf("unexpected command output")
		}
		return strings.TrimSpace(match[1]), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestVer, err := func() (string, error) {
		homepageContent, err := utils.Fetch("https://www.advancedrenamer.com/")
		if err != nil {
			return "", fmt.Errorf("can't fetch homepage: %v", err)
		}
		match := arenLatestVerRe.FindStringSubmatch(string(homepageContent))
		if match == nil {
			return "", fmt.Errorf("unexpected website content")
		}
		return strings.TrimSpace(match[1]), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	directURL := fmt.Sprintf("https://www.advancedrenamer.com/down/win/advanced_renamer_portable_%s.zip", strings.Replace(latestVer, ".", "_", -1))

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: "https://www.advancedrenamer.com/download",
	}
}
