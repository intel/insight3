package utils

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/intel-sandbox/kube-score/pkg/actions/vulns"
	"github.com/intel-sandbox/kube-score/pkg/reports"
)

func ParseRepositoryOwner(giturl string) string {
	var owner string
	gParts := strings.Split(giturl, "/")
	if len(gParts) >= 4 {
		owner = gParts[3]
	}
	return owner
}

func ParseRepositoryName(giturl string) string {
	var repo string
	repoURL := strings.TrimRight(giturl, "/")
	urlParts := strings.Split(repoURL, "/")
	repo = urlParts[len(urlParts)-1]
	return repo
}

func PrintReleaseImages(report reports.ReleaseReport) {
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

func PrintReleaseList(reports []reports.ReleaseMD) {
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

func PrintRecommendationReport(report reports.RecommendationReport) {
	fmt.Printf(strings.Repeat("*", 80))
	fmt.Printf("\nkube-score recommendation report\n")
	fmt.Printf(strings.Repeat("*", 80))

	fmt.Printf("\nRelease Measures: ")
	fmt.Printf("\n\tCurrent version: %s", report.CurrentRelease)
	fmt.Printf("\n\tLatest version: %s", report.LatestRelease)
	fmt.Printf("\n\tRecommended version: %s", report.RecommendedRelease)
	fmt.Printf("\n\tRelease lag (versions): %d", report.ReleaseLagSpace)
	fmt.Printf("\n\tRelease lag (days): %s\n", report.ReleaseLagTime)
}

func PrintVulnerabilityReport(report [][]vulns.VulnerabilityReport) {
	for i, group := range report {
		if len(group) == 0 {
			continue
		}

		fmt.Printf("\n=== Report Group %d ===\n", i+1)

		for _, r := range group {
			fmt.Printf("\nTarget: %s\n", r.Target)
			fmt.Printf("Type:   %s\n", r.Type)

			if len(r.Vulnerabilities) == 0 {
				fmt.Println("No vulnerabilities found.")
				continue
			}

			fmt.Println("Vulnerabilities:")
			for _, v := range r.Vulnerabilities {
				fmt.Printf("\n  ID:                %s\n", v.Id)
				fmt.Printf("  Package:           %s\n", v.PkgName)
				fmt.Printf("  Installed Version: %s\n", v.InstalledVersion)
				fmt.Printf("  Fixed Version:     %s\n", v.FixedVersion)
				fmt.Printf("  Title:             %s\n", v.Title)
				fmt.Printf("  Severity:          %s\n", v.Severity)
				if len(v.CWEs) > 0 {
					fmt.Printf("  CWEs:              %s\n", v.CWEs)
				}
				fmt.Printf("  Published:         %s\n", v.PublishedAt.Format(time.RFC3339))
				fmt.Printf("  Last Modified:     %s\n", v.LastModified.Format(time.RFC3339))
			}
		}
		fmt.Println()
	}
}

// Helper function to parse CreatedAt safely
func ParseTime(t string) time.Time {
	parsed, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return time.Time{} // fallback zero time if malformed
	}
	return parsed
}

// Sort and trim
func SortAndTrimReports(reports []reports.ImageReport, max int) []reports.ImageReport {
	sort.Slice(reports, func(i, j int) bool {
		return ParseTime(reports[i].CreatedAt).After(ParseTime(reports[j].CreatedAt))
	})

	if len(reports) > max {
		return reports[:max]
	}
	return reports
}
