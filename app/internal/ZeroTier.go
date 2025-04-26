package check

import (
	"fmt"
	"os/exec"
	"regexp"
	"update/utils"
)

var (
	zeroTierCurrVerRe   = regexp.MustCompile(`version ([\d\.]+)`)
	zeroTierLatestVerRe = regexp.MustCompile(`Download Latest Version \| ([\d\.]+)`)
)

func ZeroTier(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		cmdOutputBytes, err := exec.Command(execPath, "--version").Output()
		if err != nil {
			return "", err
		}
		match := zeroTierCurrVerRe.FindStringSubmatch(string(cmdOutputBytes))
		if match == nil {
			return "", err
		}
		return match[1], nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestVer, err := func() (string, error) {
		// https://www.zerotier.com/download/
		content, err := utils.Fetch("https://www.zerotier.com/download/")
		if err != nil {
			return "", fmt.Errorf("can't fetch the webpage: %v", err)
		}
		match := zeroTierLatestVerRe.FindStringSubmatch(string(content))
		if match == nil {
			return "", fmt.Errorf("unexpected website content")
		}
		return match[1], nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	directURL := "https://www.zerotier.com/download/#windows"

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: "https://www.zerotier.com/download/",
	}
}
