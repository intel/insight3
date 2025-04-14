package imageutils

import (
	"encoding/json"

	"github.com/google/go-containerregistry/pkg/crane"
)

func GetImageBuildTime(imageURL string) string {
	type imageConfig struct {
		Created string `json:"created"`
	}

	imgCfg := imageConfig{}
	c, err := crane.Config(imageURL)
	if err != nil {
		return ""
	}
	json.Unmarshal(c, &imgCfg)
	return imgCfg.Created
}

func GetImageDigest(imageURL string) string {
	d, err := crane.Digest(imageURL)
	if err != nil {
		return ""
	}
	return d
}
