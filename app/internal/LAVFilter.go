package check

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"update/utils"
)

var lavFilterCurrVerRe = regexp.MustCompile(`([\d\.]+)`)

func LAVFilter(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		changelogPath, err := getExec()
		if err != nil {
			return "", err
		}

		file, err := os.Open(changelogPath)
		if err != nil {
			return "", err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if lavFilterCurrVerRe.MatchString(line) {
				match := lavFilterCurrVerRe.FindStringSubmatch(line)
				if match != nil {
					return match[1], nil
				}
			}
		}
		err = scanner.Err()
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("can't find current version in changelog")
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestTag, latestVer, err := utils.GetGitHubLatestTag("nevcairiel/lavfilters")
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	directURL := utils.ToGitHubDirectURL(
		"nevcairiel/lavfilters",
		latestTag,
		fmt.Sprintf("LAVFilters-%s-Installer.exe", latestVer),
	)

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: fmt.Sprintf("https://github.com/lavf/lavf/releases/tag/%s", latestTag),
	}

}
