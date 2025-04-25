package ghclient

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/intel-sandbox/kube-score/pkg/common"
	"github.com/intel-sandbox/kube-score/pkg/reports"
	"github.com/pkg/errors"

	"github.com/google/go-github/github"
)

var releaseFilers = []string{
	"rc",
	"alpha",
	"beta",
}

// GetAllReleases :
func (ghcli *GHClient) GetAllReleases(ctx context.Context, repoURL string, threshold int) ([]reports.ReleaseMD, error) {
	releases := []reports.ReleaseMD{}
	repoOwner := ""
	if repoURL != "" {
		repoOwner = parseRepositoryOwner(repoURL)
	}
	repo := parseRepositoryName(repoURL)
	pageCount := 0
	for {
		pageCount++
		rlist, ghres, err := ghcli.ClientV3.Repositories.ListReleases(ctx, repoOwner, repo, &github.ListOptions{Page: pageCount, PerPage: 100})
		if err != nil || ghres.StatusCode != 200 || len(rlist) == 0 {
			if ghres.StatusCode == http.StatusForbidden {
				return releases, err
			}
			break
		}
		for _, r := range rlist {
			ignoreRelease := false
			for _, f := range releaseFilers {
				if strings.Contains(r.GetTagName(), f) {
					ignoreRelease = true
				}
			}
			if !ignoreRelease {
				releases = append(releases, reports.ReleaseMD{
					Tag:       r.GetTagName(),
					CreatedAt: r.GetCreatedAt().Time,
					CommitID:  r.GetTargetCommitish(),
				})
			}
		}
	}
	sort.Slice(releases, func(i, j int) bool {
		return releases[j].CreatedAt.Before(releases[i].CreatedAt)
	})
	if threshold != -1 && len(releases) > threshold {
		return releases[:threshold], nil
	}
	return releases, nil
}

// GetAllReleases :
func (ghcli *GHClient) GetAllReleasesGreaterThan(ctx context.Context, repoURL, currentVersion string) ([]reports.ReleaseMD, error) {
	releases := []reports.ReleaseMD{}

	rlist, err := ghcli.GetAllReleases(ctx, repoURL, -1)
	if err != nil {
		return nil, fmt.Errorf("error getting release versions")
	}
	for _, r := range rlist {
		if common.IsGreater(r.Tag, currentVersion) {
			releases = append(releases, r)
		}
	}

	sort.Slice(releases, func(i, j int) bool {
		return common.IsGreater(releases[i].Tag, releases[j].Tag)
	})

	return releases, nil
}

// GetLatestRelease :
func (ghcli *GHClient) GetLatestRelease(ctx context.Context, repoURL string) (reports.ReleaseMD, error) {
	release := reports.ReleaseMD{}
	repoOwner := ""
	if repoURL != "" {
		repoOwner = parseRepositoryOwner(repoURL)
	}
	repo := parseRepositoryName(repoURL)
	result, ghresp, err := ghcli.ClientV3.Repositories.GetLatestRelease(ctx, repoOwner, repo)
	if err != nil {
		if ghresp.StatusCode == http.StatusForbidden {
			return release, err
		}
		return release, errors.Wrapf(err, "error quering releases")
	}
	if ghresp.StatusCode != 200 {
		return release, errors.Wrapf(err, "un-expected response code %d\n", ghresp.StatusCode)
	}
	release.Tag = result.GetTagName()
	release.CreatedAt = result.GetCreatedAt().Time
	release.CommitID = result.GetTargetCommitish()
	return release, nil
}

func (ghcli *GHClient) GetReleaseTimestamp(ctx context.Context, repoURL, releaseID string) (time.Time, error) {
	ts := time.Time{}
	owner := ""
	if repoURL != "" {
		owner = parseRepositoryOwner(repoURL)
	}
	repo := parseRepositoryName(repoURL)

	release, ghresp, err := ghcli.ClientV3.Repositories.GetReleaseByTag(ctx, owner, repo, releaseID)
	if err != nil {
		fmt.Printf("error reading release time: %v", err)
		return ts, errors.Wrapf(err, "error quering releases")
	}
	if ghresp.StatusCode != 200 {
		fmt.Println(ghresp.StatusCode)
		return ts, errors.Wrapf(err, "un-expected response code %d\n", ghresp.StatusCode)
	}

	return release.GetPublishedAt().Time, nil
}

func (ghcli *GHClient) GetReleaseAsset(ctx context.Context, repoURL, releaseTag, assetName string) ([]byte, error) {
	owner := ""
	if repoURL != "" {
		owner = parseRepositoryOwner(repoURL)
	}
	repo := parseRepositoryName(repoURL)

	rlmd, _ := ghcli.GetRelease(ctx, repoURL, releaseTag)
	result, ghresp, err := ghcli.ClientV3.Repositories.ListReleaseAssets(ctx, owner, repo, rlmd.ID, &github.ListOptions{})
	if err != nil {
		if ghresp.StatusCode == http.StatusForbidden {
			return nil, err
		}
		return nil, errors.Wrapf(err, "error quering releases")
	}
	if ghresp.StatusCode != 200 {
		return nil, errors.Wrapf(err, "un-expected response code %d\n", ghresp.StatusCode)
	}

	var assetID int64
	for _, ra := range result {
		if strings.Compare(ra.GetName(), assetName) == 0 {
			assetID = ra.GetID()
			break
		}
	}

	if assetID == 0 {
		fmt.Printf("release asset not found for release [%s], assetname [%s]\n", releaseTag, assetName)
		return nil, nil
	}

	reader, redirectURL, err := ghcli.ClientV3.Repositories.DownloadReleaseAsset(ctx, owner, repo, assetID)
	if err != nil {
		if ghresp.StatusCode == http.StatusForbidden {
			return nil, err
		}
		return nil, errors.Wrapf(err, "error quering releases")
	}
	if reader != nil {
		return io.ReadAll(reader)
	} else if redirectURL != "" {
		return downloadFromURL(redirectURL)
	}

	return nil, errors.New("empty response")
}

func (ghcli *GHClient) GetRelease(ctx context.Context, repoURL, tag string) (reports.ReleaseMD, error) {
	rmd := reports.ReleaseMD{}
	repoOwner := ""
	if repoURL != "" {
		repoOwner = parseRepositoryOwner(repoURL)
	}
	repo := parseRepositoryName(repoURL)
	result, ghresp, err := ghcli.ClientV3.Repositories.GetReleaseByTag(ctx, repoOwner, repo, tag)
	if err != nil {
		if ghresp.StatusCode == http.StatusForbidden {
			return rmd, err
		}
		return rmd, errors.Wrapf(err, "error quering releases")
	}
	if ghresp.StatusCode != 200 {
		return rmd, errors.Wrapf(err, "un-expected response code %d\n", ghresp.StatusCode)
	}

	rmd.Tag = result.GetTagName()
	rmd.Name = result.GetName()
	rmd.ID = result.GetID()
	rmd.CreatedAt = result.GetCreatedAt().Time
	rmd.CommitID = result.GetTargetCommitish()

	return rmd, nil
}

// ScoreUpgradeCandidate returns a score for a candidate upgrade version.
// higher score means more recommended.
func (ghcli *GHClient) ScoreUpgradeCandidate(ctx context.Context, current, candidate string) (score float64, skip bool) {

	// skip if candidate is EOL
	if ghcli.IsEOLVersion(candidate) {
		return 0.0, true
	}

	// get vulnerabilities fixed in the candidate version
	vulns, _ := ghcli.GetFixedVulnerabilities(ctx, candidate)
	vulnScore := float64(len(vulns)) * 2.0 // higher score for more fixed vulnerabilities

	// get the changelog for the candidate version
	changelog, _ := ghcli.GetChangelog(ctx, candidate)
	changeImpact := 0.0
	if strings.Contains(changelog, "breaking") { //simple check for breaking changes
		changeImpact = 2.0
	}

	// calculate proximity score try to minimize jump distance
	versionDistance := ghcli.VersionDistance(current, candidate)
	proximityScore := 1.0 / float64(versionDistance)

	// final score
	score = vulnScore + proximityScore*1.5 - changeImpact*1.0

	return score, false
}

// GetFixedVulnerabilities fetches CVE data for a Kubernetes version from GitHub
func (ghcli *GHClient) GetFixedVulnerabilities(ctx context.Context, version string) ([]string, error) {
	// URL for GitHub's public data
	url := fmt.Sprintf("https://raw.githubusercontent.com/kubernetes/kubernetes/master/CHANGELOG/CHANGELOG-%s.md", version)

	// fetch the changelog or advisories
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching advisories from GitHub: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}

	// convert body to string
	changelogContent := string(body)

	// parse CVEs from changelog
	vulnerabilities := parseCVEsFromChangelog(changelogContent)

	return vulnerabilities, nil
}

// parseCVEsFromChangelog parses CVEs from the changelog or advisory content.
func parseCVEsFromChangelog(content string) []string {
	var vulnerabilities []string

	// Simple regex pattern to extract CVE IDs from changelog text
	// Regex for CVEs: CVE-xxxx-yyyy
	// Example changelog entry: "CVE-2021-12345: Description"
	for _, line := range strings.Split(content, "\n") {
		if strings.Contains(line, "CVE-") {
			// extract CVE ID
			cve := extractCVEID(line)
			if cve != "" {
				vulnerabilities = append(vulnerabilities, cve)
			}
		}
	}

	return vulnerabilities
}

// extractCVEID extracts the CVE ID from a line of text
func extractCVEID(line string) string {
	if strings.Contains(line, "CVE-") {
		start := strings.Index(line, "CVE-")
		if start != -1 {
			end := strings.Index(line[start:], " ")
			if end == -1 {
				end = len(line)
			}
			return line[start : start+end]
		}
	}
	return ""
}

// GetChangelog returns the changelog body for the given release version.
func (ghcli *GHClient) GetChangelog(ctx context.Context, version string) (string, error) {
	repoOwner := "kubernetes"
	repo := "kubernetes"
	release, _, err := ghcli.ClientV3.Repositories.GetReleaseByTag(ctx, repoOwner, repo, version)
	if err != nil {
		return "", err
	}
	return *release.Body, nil
}

// IsEOLVersion checks if the given version is EOL (End-Of-Life).
func (ghcli *GHClient) IsEOLVersion(version string) bool {
	eolVersions := map[string]bool{
		"v1.23": true,
		"v1.24": true,
	}
	majorMinor := version[:5]
	// add logic to fetch eol details from API
	return eolVersions[majorMinor]
}

// VersionDistance calculates the "distance" between two Kubernetes versions.
func (ghcli *GHClient) VersionDistance(current, candidate string) int {
	// Simple comparison (assuming version format v1.xx)
	currentParts := strings.Split(current, ".")
	candidateParts := strings.Split(candidate, ".")

	currentMajor, _ := strconv.Atoi(currentParts[0][1:]) // strip 'v'
	candidateMajor, _ := strconv.Atoi(candidateParts[0][1:])

	currentMinor, _ := strconv.Atoi(currentParts[1])
	candidateMinor, _ := strconv.Atoi(candidateParts[1])

	// Calculate the version distance
	majorDiff := abs(currentMajor - candidateMajor)
	minorDiff := abs(currentMinor - candidateMinor)

	return majorDiff*10 + minorDiff // weighted towards major difference
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
