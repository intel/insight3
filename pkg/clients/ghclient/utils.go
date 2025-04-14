package ghclient

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func parseRepositoryOwner(giturl string) string {
	var owner string
	gParts := strings.Split(giturl, "/")
	if len(gParts) >= 4 {
		owner = gParts[3]
	}
	return owner
}

func parseRepositoryName(giturl string) string {
	var repo string
	repoURL := strings.TrimRight(giturl, "/")
	urlParts := strings.Split(repoURL, "/")
	repo = urlParts[len(urlParts)-1]
	return repo
}

func contains(cache map[int]struct{}, key int) bool {
	if _, ok := cache[key]; ok {
		return true
	}
	return false
}

func downloadFromURL(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error reading from url")
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error reading from url")
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading from url")
	}
	return bytes, nil
}
