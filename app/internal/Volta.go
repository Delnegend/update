package check

import (
	"fmt"
	"os/exec"
	"strings"
	"update/utils"
)

func Volta(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		cmdOutputBytes, err := exec.Command(execPath, "--version").Output()
		if err != nil {
			return "", fmt.Errorf("can't run adb --version: %v", err)
		}
		return strings.TrimSpace(string(cmdOutputBytes)), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestTag, latestVer, err := utils.GetGitHubLatestTag("volta-cli/volta")
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	directURL := utils.ToGitHubDirectURL(
		"volta-cli/volta",
		latestTag,
		fmt.Sprintf("volta-%s-windows-x86_64.msi", latestVer),
	)

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: fmt.Sprintf("https://github.com/volta-cli/volta/releases/tag/%s", latestTag),
	}
}
