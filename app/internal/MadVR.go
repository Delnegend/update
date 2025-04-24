package check

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"update/utils"
)

var madVRLatestVerRe = regexp.MustCompile(`<title>madVR (\d+\.\d+\.\d+)`)

func MadVR(getExec func() (string, error)) CheckResult {
	currVer, err := func() (string, error) {
		execPath, err := getExec()
		if err != nil {
			return "", err
		}

		file, err := os.Open(execPath)
		if err != nil {
			return "", err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		if !scanner.Scan() {
			return "", fmt.Errorf("can't read the first line of changelog.txt")
		}
		return strings.TrimPrefix(strings.TrimSpace(scanner.Text()), "v"), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get current version: %v", err)}
	}

	latestVer, err := func() (string, error) {
		webpageContent, err := utils.Fetch("https://www.videohelp.com/software/madVR")
		if err != nil {
			return "", fmt.Errorf("can't fetch webpage: %v", err)
		}
		match := madVRLatestVerRe.FindStringSubmatch(string(webpageContent))
		if match == nil {
			return "", fmt.Errorf("unexpected website content")
		}
		return strings.TrimSpace(match[1]), nil
	}()
	if err != nil {
		return CheckResult{Error: fmt.Errorf("can't get latest version: %v", err)}
	}

	directURL := fmt.Sprintf("https://www.videohelp.com/download/madVR%s.zip", strings.Replace(latestVer, ".", "", -1))

	return CheckResult{
		CurrentVersion:  currVer,
		LatestVersion:   latestVer,
		DirectURL:       directURL,
		DirectURLAlive:  utils.IsURLAlive(directURL, currVer == latestVer),
		DownloadPageURL: "https://www.videohelp.com/software/madVR",
	}
}
