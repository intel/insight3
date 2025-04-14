package reports

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
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
	ReleaseLagTime     string `json:"releaseLagTime"`
	ReleaseLagSpace    int    `json:"releaseLagSpace"`
}

func PrintReleaseImages(report ReleaseReport) {
	fmt.Printf(strings.Repeat("*", 80))
	fmt.Printf("\nkube-score release report\n")
	fmt.Printf(strings.Repeat("*", 80))
	fmt.Println()
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, fmt.Sprintf("%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s", "ImageURL", "Tag", "Digest", "BuildTime", "Signed", "Vulnerabilities"))
	fmt.Fprintln(w, fmt.Sprintf("%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s", strings.Repeat("-", 5), strings.Repeat("-", 5), strings.Repeat("-", 5), strings.Repeat("-", 5), strings.Repeat("-", 5), strings.Repeat("-", 5)))

	for _, i := range report.Images {
		parts := strings.Split(i.URL, ":")

		if len(parts) != 2 {
			return
		}
		url := parts[0]
		tag := parts[1]
		dRunes := []rune(i.Digest)
		digestSmall := string(dRunes[0:14])
		vuln := fmt.Sprintf("C[%d],H[%d],M[%d],L[%d]", i.Vulnerabilities.Summary.Critical, i.Vulnerabilities.Summary.High, i.Vulnerabilities.Summary.Medium, i.Vulnerabilities.Summary.Low)
		fmt.Fprintln(w, fmt.Sprintf("%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s", url, tag, digestSmall, i.CreatedAt, "yes", vuln))
	}
	w.Flush()
}

func PrintReleaseList(reports []ReleaseMD) {
	fmt.Printf(strings.Repeat("*", 80))
	fmt.Printf("\nkube-score release report\n")
	fmt.Printf(strings.Repeat("*", 80))
	fmt.Println()
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, fmt.Sprintf("%s\t\t%s", "ReleaseTag", "ReleaseTime"))
	fmt.Fprintln(w, fmt.Sprintf("%s\t\t%s", strings.Repeat("-", 5), strings.Repeat("-", 5)))

	for _, r := range reports {
		fmt.Fprintln(w, fmt.Sprintf("%s\t\t%s", r.Tag, r.CreatedAt))
	}
	w.Flush()
}

func PrintRecommendationReport(report RecommendationReport) {
	fmt.Printf(strings.Repeat("*", 80))
	fmt.Printf("\nkube-score recommendation report\n")
	fmt.Printf(strings.Repeat("*", 80))

	fmt.Printf("\nRelease Measures: ")
	fmt.Printf("\n\tCurrent version: %s", report.CurrentRelease)
	fmt.Printf("\n\tLatest version: %s", report.RecommendedRelease)
	fmt.Printf("\n\tRecommended version: %s", report.RecommendedRelease)
	fmt.Printf("\n\tRelease lag (versions): %d", report.ReleaseLagSpace)
	fmt.Printf("\n\tRelease lag (days): %s\n", report.ReleaseLagTime)
}
