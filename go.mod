module github.com/intel-sandbox/kube-score

go 1.23.0

toolchain go1.23.8

require (
	github.com/Masterminds/semver v1.5.0
	github.com/go-git/go-git/v5 v5.16.0
	github.com/google/go-containerregistry v0.12.1
	github.com/google/go-github v17.0.0+incompatible
	github.com/hako/durafmt v0.0.0-20210608085754-5c1018a4e16b
	github.com/pkg/errors v0.9.1
	github.com/shurcooL/githubv4 v0.0.0-20221203213311-70889c5dac07
	github.com/spf13/cobra v1.8.1
	go.uber.org/zap v1.27.0
	golang.org/x/oauth2 v0.27.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/kubernetes v1.32.6
)

require (
	dario.cat/mergo v1.0.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/ProtonMail/go-crypto v1.1.6 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cloudflare/circl v1.6.1 // indirect
	github.com/cyphar/filepath-securejoin v0.4.1 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.6.2 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pjbgf/sha1cd v0.3.2 // indirect
	github.com/prometheus/client_golang v1.19.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/sergi/go-diff v1.3.2-0.20230802210424-5b0b94c5c0d3 // indirect
	github.com/skeema/knownhosts v1.3.1 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	go.opentelemetry.io/otel v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.12.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/docker/cli v20.10.20+incompatible // indirect
	github.com/docker/distribution v2.8.2+incompatible // indirect
	github.com/docker/docker v25.0.6+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.7.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-redis/redis/v8 v8.11.5
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.15.11 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc2 // indirect
	github.com/shurcooL/graphql v0.0.0-20220606043923-3cf50f8a0a29 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/vbatts/tar-split v0.11.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.32.4 // indirect
	k8s.io/apimachinery v0.32.4 // indirect
	k8s.io/client-go v0.32.4 // indirect
	k8s.io/cluster-bootstrap v0.0.0 // indirect
	k8s.io/component-base v0.0.0 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/utils v0.0.0-20241104100929-3ea5e8cea738 // indirect
	sigs.k8s.io/json v0.0.0-20241010143419-9aa6b5e7a4b3 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.2 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

replace k8s.io/api => k8s.io/api v0.32.4

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.32.4

replace k8s.io/apimachinery => k8s.io/apimachinery v0.32.4

replace k8s.io/apiserver => k8s.io/apiserver v0.32.4

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.32.4

replace k8s.io/client-go => k8s.io/client-go v0.32.4

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.32.4

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.32.4

replace k8s.io/code-generator => k8s.io/code-generator v0.32.4

replace k8s.io/component-base => k8s.io/component-base v0.32.4

replace k8s.io/component-helpers => k8s.io/component-helpers v0.32.4

replace k8s.io/controller-manager => k8s.io/controller-manager v0.32.4

replace k8s.io/cri-api => k8s.io/cri-api v0.32.4

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.32.4

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.32.4

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.32.4

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.32.4

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.32.4

replace k8s.io/kubectl => k8s.io/kubectl v0.32.4

replace k8s.io/kubelet => k8s.io/kubelet v0.32.4

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.32.4

replace k8s.io/metrics => k8s.io/metrics v0.32.4

replace k8s.io/mount-utils => k8s.io/mount-utils v0.32.4

replace k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.32.4

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.32.4

replace k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.32.4

replace k8s.io/sample-controller => k8s.io/sample-controller v0.32.4
