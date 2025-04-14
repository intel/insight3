package config

import (
	"fmt"
	"io/ioutil"

	"github.com/intel-sandbox/kube-score/pkg/common"
	"gopkg.in/yaml.v2"
)

func ConfigParser(cfgFilepath *string) *common.RunConfig {

	cfg := common.RunConfig{}
	yamlBuf, err := ioutil.ReadFile(*cfgFilepath)
	if err != nil {
		fmt.Printf("error reading control file: %s: %v", *cfgFilepath, err)
		return nil
	}
	if err := yaml.Unmarshal(yamlBuf, &cfg); err != nil {
		fmt.Printf("error parsing control file: %s: %v", *cfgFilepath, err)
		return nil
	}
	return &cfg
}
