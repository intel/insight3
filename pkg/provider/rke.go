package provider

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/intel-sandbox/kube-score/pkg/actions/imageutils"
	"github.com/intel-sandbox/kube-score/pkg/clients/db"
	"github.com/intel-sandbox/kube-score/pkg/clients/ghclient"
	"github.com/intel-sandbox/kube-score/pkg/reports"
)

type RkeProvider struct {
	RepoURL  string
	GHClient *ghclient.GHClient
	DBClient *db.RedisClient
}

const (
	releaseAssetImageList = "rke2-images-all.linux-amd64.txt"
)

func (rke RkeProvider) GetReleases(ctx context.Context) ([]reports.ReleaseMD, error) {
	return rke.GHClient.GetAllReleases(ctx, rke.RepoURL, 20)
}

func (rke RkeProvider) GetReleaseMeta(ctx context.Context, version string) (reports.ReleaseMD, error) {
	return rke.GHClient.GetRelease(ctx, rke.RepoURL, version)
}

func (rke RkeProvider) GetReleaseAssets(ctx context.Context, version, name string) ([]reports.ReleaseAsset, error) {

	return nil, nil
}

func (rke RkeProvider) GetReleaseImages(ctx context.Context, version string) ([]reports.ImageReport, error) {
	imagesBuf, err := rke.GHClient.GetReleaseAsset(ctx, rke.RepoURL, version, releaseAssetImageList)
	if err != nil {
		fmt.Printf("error reading release asset for repo [%s], release [%s]\n", rke.RepoURL, releaseAssetImageList)
		return nil, errors.New("failed to read release images")
	}
	reader := bytes.NewReader(imagesBuf)
	fileScanner := bufio.NewScanner(reader)

	fileScanner.Split(bufio.ScanLines)
	imgs := []reports.ImageReport{}
	for fileScanner.Scan() {
		imgs = append(imgs, reports.ImageReport{
			URL:       fileScanner.Text(),
			CreatedAt: imageutils.GetImageBuildTime(fileScanner.Text()),
			Digest:    imageutils.GetImageDigest(fileScanner.Text()),
		})
	}

	return imgs, nil
}
