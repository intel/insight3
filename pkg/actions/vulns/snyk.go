package vulns

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"html/template"
// 	"os"
// 	"os/exec"
// 	"strings"
// 	"syscall"

// 	"github.com/intel-sandbox/kube-score/pkg/common"
// 	"github.com/intel-sandbox/kube-score/pkg/reports"
// )

// const (
// 	snykToken = "SNYK_TOKEN"
// 	snykApi   = "SNYK_API"

// 	scanContainerCmd = "snyk container test {{.URL}} --app-vulns --json"
// )

// type SnykScanner struct {
// }

// type SnykReport struct {
// 	Path            string         `json:"path"`
// 	Summary         string         `json:"summary"`
// 	ProjectName     string         `json:"projectName"`
// 	UniqueCount     int            `json:"uniqueCount"`
// 	Platform        string         `json:"platform"`
// 	Vulnerabilities []Vulnerbility `json:"vulnerabilities"`
// 	Applications    []struct {
// 		Path              string         `json:"path"`
// 		UnkownVersion     bool           `json:"hasUnknownVersions"`
// 		DisplayTargetFile string         `json:"displayTargetFile"`
// 		ProjectName       string         `json:"projectName"`
// 		TargetFile        string         `json:"targetFile"`
// 		UniqueCount       int            `json:"uniqueCount"`
// 		DependencyCount   int            `json:"dependencyCount"`
// 		Vulnerabilities   []Vulnerbility `json:"vulnerabilities"`
// 	} `json:"applications"`
// }

// type Vulnerbility struct {
// 	Id               string  `json:"id"`
// 	CVSSv3           string  `json:"CVSSv3"`
// 	CreationTime     string  `json:"creationTime"`
// 	CVSSScore        float64 `json:"cvssScore"`
// 	DisclosureTime   string  `json:"disclosureTime"`
// 	ModificationTime string  `json:"modificationTime"`
// 	PublicationTime  string  `json:"publicationTime"`
// 	Exploit          string  `json:"exploit"`
// 	Severity         string  `json:"severity"`
// 	Title            string  `json:"title"`
// 	Semver           struct {
// 		Vulnerable []string `json:"vulnerable"`
// 	} `json:"semver"`
// 	Identifiers struct {
// 		CWE []string `json:"CWE"`
// 		CVE []string `json:"CVE"`
// 	} `json:"identifiers"`
// 	Language   string   `json:"language"`
// 	ModuleName string   `json:"moduleName"`
// 	From       []string `json:"from"`
// 	Name       string   `json:"name"`
// 	Version    string   `json:"version"`
// }

// func (snykCli SnykScanner) Init(opts common.ScannerConfig) error {
// 	// setup environment
// 	os.Setenv(snykApi, opts.Snyk.Endpoint)
// 	os.Setenv(snykToken, opts.Snyk.AuthToken)

// 	return nil
// }

// func (snykCli SnykScanner) ScanImage(ctx context.Context, imageURL string) (reports.VulnerabilityData, error) {
// 	vulns := reports.VulnerabilityData{}

// 	cmdTmpl, err := template.New("snykScanCmd").Parse(scanContainerCmd)
// 	if err != nil {
// 		fmt.Printf("error formating snyk scan cmd: %v", err)
// 	}

// 	type input struct {
// 		URL string
// 	}

// 	var cmdBuf bytes.Buffer
// 	cmdTmpl.Execute(&cmdBuf, input{URL: imageURL})
// 	report, err := execCmd(cmdBuf.String())
// 	if err != nil {
// 		fmt.Printf("failed to scan image: %s, err: %v\n", imageURL, err)
// 		return vulns, fmt.Errorf("failed to scan image: %s", imageURL)
// 	}
// 	snykReport := SnykReport{}
// 	json.Unmarshal(report, &snykReport)
// 	updateVulnerabilityStatus(&vulns, &snykReport)

// 	return vulns, nil
// }

// func execCmd(cmdStr string) ([]byte, error) {
// 	cmdS := strings.Split(cmdStr, " ")
// 	cmd := exec.Command(cmdS[0], cmdS[1:]...)
// 	var outb, errb bytes.Buffer
// 	cmd.Stdout = &outb
// 	cmd.Stderr = &errb
// 	var waitStatus syscall.WaitStatus
// 	exitCode := 0
// 	if err := cmd.Run(); err != nil {
// 		// Did the command fail because of an unsuccessful exit code
// 		if exitError, ok := err.(*exec.ExitError); ok {
// 			waitStatus = exitError.Sys().(syscall.WaitStatus)
// 			exitCode = waitStatus.ExitStatus()
// 		}
// 	} else {
// 		// Command was successful
// 		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
// 		exitCode = waitStatus.ExitStatus()
// 	}
// 	if exitCode > 1 {
// 		return nil, fmt.Errorf("failed to scan image")
// 	}

// 	return outb.Bytes(), nil
// }

// func updateVulnerabilityStatus(vulns *reports.VulnerabilityData, snykReport *SnykReport) {
// 	total := 0
// 	for _, osV := range snykReport.Vulnerabilities {
// 		switch osV.Severity {
// 		case "critical":
// 			vulns.Summary.Critical++
// 			total++
// 		case "high":
// 			vulns.Summary.High++
// 			total++
// 		case "medium":
// 			vulns.Summary.Medium++
// 			total++
// 		case "low":
// 			vulns.Summary.Low++
// 			total++
// 		}
// 	}

// 	for _, app := range snykReport.Applications {
// 		for _, appV := range app.Vulnerabilities {
// 			switch appV.Severity {
// 			case "critical":
// 				vulns.Summary.Critical++
// 				total++
// 			case "high":
// 				vulns.Summary.High++
// 				total++
// 			case "medium":
// 				vulns.Summary.Medium++
// 				total++
// 			case "low":
// 				vulns.Summary.Low++
// 				total++
// 			}
// 		}
// 	}
// 	vulns.Summary.Total = total
// }
