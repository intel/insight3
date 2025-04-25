// INTEL CONFIDENTIAL
// Copyright (C) 2023 Intel Corporation
package vulns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/intel-sandbox/kube-score/pkg/reports"
)

const (
	// Define the command template
	scanContainerCmd = "{{.Trivy}} image -f json -o {{.OutputFile}} {{.URL}}"
	scanRepoCmd      = "{{.Trivy}} fs -f json -o {{.OutputFile}} {{.RepoPath}}"
)

type Vulnerability struct {
	Id               string    `json:"VulnerabilityID"`
	PkgName          string    `json:"PkgName"`
	InstalledVersion string    `json:"InstalledVersion"`
	Title            string    `json:"Title"`
	Severity         string    `json:"Severity"`
	CWEs             []string  `json:"CweIDs"`
	PublishedAt      time.Time `json:"PublishedDate"`
	LastModified     time.Time `json:"LastModifiedDate"`
	FixedVersion     string    `json:"FixedVersion"`
}

type VulnerabilityReport struct {
	Target          string          `json:"Target"`
	Type            string          `json:"Type"`
	Vulnerabilities []Vulnerability `json:"Vulnerabilities"`
}

type TrivyResult struct {
	SchemaVersion int                   `json:"SchemaVersion"`
	Results       []VulnerabilityReport `json:"Results"`
}
type TrivyScanner struct {
}

func (scanner *TrivyScanner) ScanImage(ctx context.Context, imageURL string) (*reports.VulnerabilityData, []VulnerabilityReport, error) {
	vulnsummary := reports.VulnerabilityData{}

	os.Setenv("https_proxy", "")
	os.Setenv("HTTPS_PROXY", "")

	// check for presence of trivy
	path, err := exec.Command("which", "trivy").Output()
	if err != nil {
		fmt.Println("failed to locate trivy binary:", err)
		return nil, nil, err
	}
	trivyPath := strings.TrimSpace(string(path))

	// grant permission to path
	os.Chmod(trivyPath, 0755)
	cmdTmpl, err := template.New("trivyScanCmd").Parse(scanContainerCmd)
	if err != nil {
		fmt.Printf("error formating trivy scan cmd: %v", err)
	}

	type input struct {
		URL        string
		OutputFile string
		Trivy      string
	}
	// create output directory to store trivy results
	outJson, err := os.CreateTemp(os.Getenv("HOME"), "")
	if err != nil {
		fmt.Println("failed to create temp directory")
		return nil, nil, err
	}
	defer outJson.Close()
	defer os.RemoveAll(outJson.Name())

	var cmdBuf bytes.Buffer
	err = cmdTmpl.Execute(&cmdBuf, input{
		URL:        imageURL,
		OutputFile: outJson.Name(),
		Trivy:      trivyPath,
	})
	if err != nil {
		fmt.Errorf("\nerror executing trivy command template: %w", err)
	}

	fmt.Printf("Executing Trivy command: %s\n", cmdBuf.String())
	if _, err = execCmd(cmdBuf.String()); err != nil {
		return nil, nil, fmt.Errorf("failed to scan image %s: %w", imageURL, err)
	}

	// read stored trivy results
	vBuf, err := os.ReadFile(outJson.Name())
	if err != nil {
		return nil, nil, fmt.Errorf("error reading Trivy output: %w", err)
	}

	trivyRes := TrivyResult{}
	if err := json.Unmarshal(vBuf, &trivyRes); err != nil {
		fmt.Println(err, "error parsing json result")
		return nil, nil, err
	}
	// format results
	updateVulnerabilityStatus(&vulnsummary, trivyRes.Results)

	return &vulnsummary, trivyRes.Results, nil
}

func (scanner *TrivyScanner) ScanRepo(ctx context.Context, repoURL string) (*reports.VulnerabilityData, []VulnerabilityReport, error) {
	vulnsummary := reports.VulnerabilityData{}
	os.Setenv("https_proxy", "")
	os.Setenv("HTTPS_PROXY", "")

	//extract component and version from url
	cname, version, err := parseGitHubReleaseURL(repoURL)
	if err != nil {
		fmt.Println("Error parsing GitHub release URL:", err)
	} else {
		fmt.Printf("Component name: %s version: %s\n", cname, version)
	}

	// check for presence of trivy
	path, err := exec.Command("which", "trivy").Output()
	if err != nil {
		fmt.Println("failed to locate trivy binary:", err)
		return nil, nil, err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("failed to get user home directory: %s\n", err)
		return nil, nil, err
	}

	scanDir := filepath.Join(homeDir, "trivy-scan-workdir", fmt.Sprintf("repo-%d", time.Now().UnixNano()))
	if err := os.MkdirAll(scanDir, 0755); err != nil {
		fmt.Printf("failed to create scan directory: %s\n", err)
		return nil, nil, err
	}
	defer os.RemoveAll(scanDir) // Cleanup temp directory after execution

	base := strings.Split(repoURL, "/releases/tag/")[0]
	repoUrl := base + ".git"
	// Clone the GitHub repo into scanDir
	_, err = git.PlainClone(scanDir, false, &git.CloneOptions{
		URL:           repoUrl,
		ReferenceName: plumbing.ReferenceName("refs/tags/" + version), // for tag
		SingleBranch:  true,
	})
	if err != nil {
		fmt.Printf("failed to clone repository: %s\n", err)
		return nil, nil, err
	}

	trivyPath := strings.TrimSpace(string(path))
	// grant permission to path
	os.Chmod(trivyPath, 0755)
	cmdTmpl, err := template.New("trivyScanCmd").Parse(scanRepoCmd)
	if err != nil {
		fmt.Printf("error formating trivy scan cmd: %v", err)
	}
	type input struct {
		RepoPath   string
		OutputFile string
		Trivy      string
	}

	// Create temporary output file for scan results
	outJson, err := os.CreateTemp(os.Getenv("HOME"), "")
	if err != nil {
		fmt.Printf("error creating temp file: %s\n", err)
		return nil, nil, err
	}

	defer outJson.Close()
	defer os.RemoveAll(outJson.Name())

	// Render the command
	var cmdBuf bytes.Buffer
	err = cmdTmpl.Execute(&cmdBuf, input{
		RepoPath:   scanDir,
		OutputFile: outJson.Name(),
		Trivy:      trivyPath,
	})
	if err != nil {
		fmt.Printf("failed to render Trivy scan command")
	}
	//fmt.Printf("Executing Trivy scan: %s\n", cmdBuf.String())

	// Execute the scan command
	_, err = execCmd(cmdBuf.String())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to scan repository: %s", repoURL)
	}

	// read stored trivy results
	vBuf, err := os.ReadFile(outJson.Name())
	if err != nil {
		return nil, nil, fmt.Errorf("error reading Trivy output: %w", err)
	}
	trivyRes := TrivyResult{}
	if err := json.Unmarshal(vBuf, &trivyRes); err != nil {
		fmt.Println(err, "error parsing json result")
		return nil, nil, err
	}

	// format results
	updateVulnerabilityStatus(&vulnsummary, trivyRes.Results)

	return &vulnsummary, trivyRes.Results, nil
}

func execCmd(cmdStr string) ([]byte, error) {
	cmdS := strings.Split(cmdStr, " ")
	cmd := exec.Command(cmdS[0], cmdS[1:]...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	var waitStatus syscall.WaitStatus
	exitCode := 0
	if err := cmd.Run(); err != nil {
		// Did the command fail because of an unsuccessful exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			exitCode = waitStatus.ExitStatus()
		}
	} else {
		// Command was successful
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = waitStatus.ExitStatus()
	}
	if exitCode > 1 {
		return nil, fmt.Errorf("failed to scan image")
	}

	return outb.Bytes(), nil
}

// helper func to count number of vuln by severity
func updateVulnerabilityStatus(summary *reports.VulnerabilityData, results []VulnerabilityReport) {
	for _, report := range results {
		for _, vuln := range report.Vulnerabilities {
			switch strings.ToLower(vuln.Severity) {
			case "critical":
				summary.Summary.Critical++
			case "high":
				summary.Summary.High++
			case "medium":
				summary.Summary.Medium++
			case "low":
				summary.Summary.Low++
			}
			summary.Summary.Total++
		}
	}
}

func parseGitHubReleaseURL(url string) (string, string, error) {
	parts := strings.Split(url, "/")
	if len(parts) < 7 {
		return "", "", fmt.Errorf("invalid GitHub release URL: %s", url)
	}

	component := parts[4]          // e.g. "runc"
	version := parts[len(parts)-1] // e.g. "v1.2.6"

	return component, version, nil
}
