package check

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"update/utils"
)

var platformToolsLatestVerRe = regexp.MustCompile(`>Revisions<(.|\s)+?data-text="([\d\.]+)`)
var platformToolsCurrVerRe = regexp.MustCompile(`Version (\d+\.\d+\.\d+)`)

func PlatformTools(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		cmdOutputBytes, err := exec.Command(execPath, "--version").Output()
		if err != nil {
			return "", fmt.Errorf("can't run adb --version: %v", err)
		}
		match := platformToolsCurrVerRe.FindStringSubmatch(string(cmdOutputBytes))
		if match == nil {
			return "", fmt.Errorf("unexpected command output")
		}
		return strings.TrimSpace(match[1]), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestVer, err := func() (string, error) {
		webpageContent, err := utils.Fetch("https://developer.android.com/tools/releases/platform-tools")
		if err != nil {
			return "", fmt.Errorf("can't fetch webpage: %v", err)
		}
		match := platformToolsLatestVerRe.FindStringSubmatch(string(webpageContent))
		if match == nil {
			return "", fmt.Errorf("unexpected website content")
		}
		return strings.TrimSpace(match[2]), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	directURL := "https://dl.google.com/android/repository/platform-tools-latest-windows.zip"

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: "https://developer.android.com/tools/releases/platform-tools",
	}
}
