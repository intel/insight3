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
	"strings"
	"syscall"
	"time"

	"github.com/intel-sandbox/kube-score/pkg/reports"
)

const (
	// Define the command template
	scanContainerCmd = "/opt/homebrew/bin/trivy image -f json -o {{.OutputFile}} {{.URL}}"
)

type VulnerabilityReport struct {
	Target          string `json:"Target"`
	Type            string `json:"Type"`
	Vulnerabilities []struct {
		Id               string    `json:"VulnerabilityID"`
		PkgName          string    `json:"PkgName"`
		InstalledVersion string    `json:"InstalledVersion"`
		Title            string    `json:"Title"`
		Severity         string    `json:"Severity"`
		CWEs             []string  `json:"CweIDs"`
		PublishedAt      time.Time `json:"PublishedDate"`
		LastModified     time.Time `json:"LastModifiedDate"`
	} `json:"Vulnerabilities"`
}

type TrivyScanner struct {
}

func (scanner *TrivyScanner) ScanImage(ctx context.Context, imageURL string) (*reports.VulnerabilityData, error) {
	vulnsummary := reports.VulnerabilityData{}
	vulns := []VulnerabilityReport{}

	os.Setenv("https_proxy", "")
	os.Setenv("HTTPS_PROXY", "")

	os.Chmod("/opt/homebrew/bin/trivy", 0755)
	cmdTmpl, err := template.New("trivyScanCmd").Parse(scanContainerCmd)
	if err != nil {
		fmt.Printf("error formating trivy scan cmd: %v", err)
	}

	type input struct {
		URL        string
		OutputFile string
	}

	outJson, err := os.CreateTemp(os.TempDir(), "")
	if err != nil {
		return nil, err
	}
	defer outJson.Close()
	defer os.RemoveAll(outJson.Name())

	var cmdBuf bytes.Buffer
	cmdTmpl.Execute(&cmdBuf, input{URL: imageURL, OutputFile: outJson.Name()})
	fmt.Printf("execute command", "cmd", cmdBuf.String())

	_, err = execCmd(cmdBuf.String())
	if err != nil {
		fmt.Printf("failed to scan image: %s, err: %v\n", imageURL, err)
		return nil, fmt.Errorf("failed to scan image: %s", imageURL)
	}

	vBuf, _ := os.ReadFile(outJson.Name())
	if err := json.Unmarshal(vBuf, &vulns); err != nil {
		fmt.Println(err, "error parsing json result")
		return nil, err
	}

	updateVulnerabilityStatus(&vulnsummary, vulns)
	return &vulnsummary, nil
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

func updateVulnerabilityStatus(vulns *reports.VulnerabilityData, trivy []VulnerabilityReport) {
	total := 0
	for _, osV := range trivy {
		for _, v := range osV.Vulnerabilities {
			switch v.Severity {
			case "critical":
				vulns.Summary.Critical++
				total++
			case "high":
				vulns.Summary.High++
				total++
			case "medium":
				vulns.Summary.Medium++
				total++
			case "low":
				vulns.Summary.Low++
				total++
			}
		}
	}
}
