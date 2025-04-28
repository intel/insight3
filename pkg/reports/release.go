package reports

import (
	"time"
)

// ReleaseMD :
type ReleaseMD struct {
	Tag       string    `json:"release_tag"`
	CreatedAt time.Time `json:"created_at"`
	CommitID  string    `json:"commit_id"`
	Name      string    `json:"name"`
	ID        int64     `json:"id"`
}

type ReleaseReport struct {
	ReleaseTag  string        `json:"releaseTag"`
	ReleaseTime time.Time     `json:"releaseTime"`
	Images      []ImageReport `json:"images"`
}

type ImageReport struct {
	URL             string            `json:"url"`
	Digest          string            `json:"digest"`
	CreatedAt       string            `json:"createdAt"`
	Vulnerabilities VulnerabilityData `json:"vulnerabilities"`
}

type VulnerabilityData struct {
	Summary VulnerabilitySummary `json:"vulnerabilities"`
}

type VulnerabilitySummary struct {
	Total    int `json:"total"`
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
}

type ReleaseAsset struct {
	ReleaseID   string
	Name        string
	Type        string
	DownloadURL string
}

type RecommendationReport struct {
	CurrentRelease     string `json:"currentRelease"`
	RecommendedRelease string `json:"recommendedRelease"`
	LatestRelease      string `json:"latestRelease"`
	ReleaseLagTime     string `json:"releaseLagTime"`
	ReleaseLagSpace    int    `json:"releaseLagSpace"`
}

