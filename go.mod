module github.com/intel-sandbox/kube-score

go 1.19

require (
	github.com/Masterminds/semver v1.5.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/google/go-containerregistry v0.12.1
	github.com/google/go-github v17.0.0+incompatible
	github.com/hako/durafmt v0.0.0-20210608085754-5c1018a4e16b
	github.com/pkg/errors v0.9.1
	github.com/shurcooL/githubv4 v0.0.0-20221203213311-70889c5dac07
	github.com/spf13/cobra v1.6.1
	go.uber.org/zap v1.24.0
	golang.org/x/oauth2 v0.3.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/kubernetes v1.25.4
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.12.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/docker/cli v20.10.20+incompatible // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v25.0.6+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.7.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.15.11 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/onsi/gomega v1.24.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc2 // indirect
	github.com/shurcooL/graphql v0.0.0-20220606043923-3cf50f8a0a29 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	github.com/vbatts/tar-split v0.11.2 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.25.4 // indirect
	k8s.io/apimachinery v0.25.4 // indirect
	k8s.io/client-go v0.25.4 // indirect
	k8s.io/cluster-bootstrap v0.0.0 // indirect
	k8s.io/component-base v0.0.0 // indirect
	k8s.io/klog/v2 v2.70.1 // indirect
	k8s.io/utils v0.0.0-20220728103510-ee6ede2d64ed // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace k8s.io/api => k8s.io/api v0.25.4

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.25.4

replace k8s.io/apimachinery => k8s.io/apimachinery v0.25.5-rc.0

replace k8s.io/apiserver => k8s.io/apiserver v0.25.4

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.25.4

replace k8s.io/client-go => k8s.io/client-go v0.25.4

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.25.4

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.25.4

replace k8s.io/code-generator => k8s.io/code-generator v0.25.5-rc.0

replace k8s.io/component-base => k8s.io/component-base v0.25.4

replace k8s.io/component-helpers => k8s.io/component-helpers v0.25.4

replace k8s.io/controller-manager => k8s.io/controller-manager v0.25.4

replace k8s.io/cri-api => k8s.io/cri-api v0.25.5-rc.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.25.4

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.25.4

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.25.4

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.25.4

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.25.4

replace k8s.io/kubectl => k8s.io/kubectl v0.25.4

replace k8s.io/kubelet => k8s.io/kubelet v0.25.4

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.25.4

replace k8s.io/metrics => k8s.io/metrics v0.25.4

replace k8s.io/mount-utils => k8s.io/mount-utils v0.25.5-rc.0

replace k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.25.4

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.25.4

replace k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.25.4

replace k8s.io/sample-controller => k8s.io/sample-controller v0.25.4
