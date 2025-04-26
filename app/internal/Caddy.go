package check

import (
	"fmt"
	"os/exec"
	"strings"
	"update/utils"
)

func Caddy(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		currentVersionBytes, err := exec.Command(execPath, "version").Output()
		if err != nil {
			return "", fmt.Errorf("unexpected command output")
		}
		currentVersionSlice := strings.SplitN(string(currentVersionBytes), " ", 2)
		if len(currentVersionSlice) == 2 {
			return currentVersionSlice[0], nil
		}
		return string(currentVersionBytes), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestTag, latestVer, err := utils.GetGitHubLatestTag("caddyserver/caddy")
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	directURL := utils.ToGitHubDirectURL(
		"caddyserver/caddy",
		latestTag,
		fmt.Sprintf("caddy_%s_windows_amd64.zip", latestVer),
	)

	return CheckResult{
		CurrentVersion:  strings.TrimPrefix(currVer, "v"),
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: fmt.Sprintf("https://github.com/caddyserver/caddy/releases/tag/%s", latestTag),
	}
}
