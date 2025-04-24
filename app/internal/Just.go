package check

import (
	"fmt"
	"os/exec"
	"strings"
	"update/utils"
)

func Just(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		currVerBytes, err := exec.Command(execPath, "--version").Output()
		if err != nil {
			return "", err
		}
		currVerString := string(currVerBytes)
		currVerSlice := strings.SplitN(currVerString, " ", 2)
		if len(currVerSlice) == 2 {
			return currVerSlice[0], nil
		}
		return "", fmt.Errorf("unexpected command output")
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestTag, latestVer, err := utils.GetGitHubLatestTag("casey/just")
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	directURL := utils.ToGitHubDirectURL(
		"casey/just",
		latestTag,
		fmt.Sprintf("just-%s-x86_64-pc-windows-msvc.zip", latestVer),
	)

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: fmt.Sprintf("https://github.com/casey/just/releases/tag/%s", latestTag),
		Error:           nil,
	}
}
