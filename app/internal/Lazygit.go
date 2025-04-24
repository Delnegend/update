package check

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"update/utils"
)

var lazygitCurrVerRe = regexp.MustCompile(`version=([\d\.]+)`)

func Lazygit(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		currVerBytes, err := exec.Command(execPath, "--version").Output()
		if err != nil {
			return "", err
		}
		match := lazygitCurrVerRe.FindStringSubmatch(string(currVerBytes))
		if match == nil {
			return "", fmt.Errorf("unexpected command output")
		}
		return strings.TrimSpace(match[1]), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestTag, latestVer, err := utils.GetGitHubLatestTag("jesseduffield/lazygit")
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	directURL := utils.ToGitHubDirectURL(
		"jesseduffield/lazygit",
		latestTag,
		fmt.Sprintf("lazygit_%s_Windows_x86_64.zip", latestVer),
	)
	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: fmt.Sprintf("https://github.com/jesseduffield/lazygit/releases/tag/%s", latestTag),
		Error:           nil,
	}
}
