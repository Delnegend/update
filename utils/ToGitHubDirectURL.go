package utils

import "fmt"

func ToGitHubDirectURL(repo, tag, fileName string) string {
	return fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", repo, tag, fileName)
}
