# Insights3 (kube-score)

Insights for kubernetes (or kube-score) is a command line tool for gathering vulnerabilities details for different Kubernetes release versions. This command line tool builds on top of existing vulnerability scanning solutions to yield a more structured release vulnerability report. Intention is for tool to aid in discovery of CVEs between releases and for informing component upgrade decisions.​


## Getting Started

### Pre-requisites

Install go 1.19+ on your dev environment

### Local Dev Setup

1. Clone this repository and run `make` inside it. It should build an executable binary under `bin/`

```
$ git clone https://github.com/intel/insight3.git

$ cd insight3

$ make
```

### Configurations
Create a file called `.kube_score_use.yaml` using `.kube_score.yaml` as reference and update with your Github PAT. It is required to avoid rate-limit throttling when calling Github APIs.
(Currently, we don't need any other configuration in this file)
![image](https://github.com/user-attachments/assets/9300bc72-b84b-4940-9877-b887b03ce6f5)


### Testing

1. Listing all recent k8s releases (capped at top-20)
```
$ ./bin/kube-score release --list
********************************************************************************
kube-score release report
********************************************************************************
ReleaseTag  ReleaseTime
-----       -----
v1.26.0     2022-12-08 19:51:43 +0000 UTC
v1.22.17    2022-12-08 11:41:04 +0000 UTC
v1.23.15    2022-12-08 10:41:43 +0000 UTC
v1.25.5     2022-12-08 10:06:36 +0000 UTC
v1.24.9     2022-12-08 10:06:34 +0000 UTC
v1.23.14    2022-11-09 13:31:23 +0000 UTC
v1.22.16    2022-11-09 13:30:52 +0000 UTC
v1.24.8     2022-11-09 13:30:14 +0000 UTC
v1.25.4     2022-11-09 13:28:30 +0000 UTC
v1.23.13    2022-10-12 10:49:33 +0000 UTC
v1.24.7     2022-10-12 10:48:53 +0000 UTC
v1.25.3     2022-10-12 10:47:25 +0000 UTC
v1.25.2     2022-09-21 14:25:45 +0000 UTC
v1.24.6     2022-09-21 13:10:31 +0000 UTC
v1.23.12    2022-09-21 12:11:42 +0000 UTC
v1.22.15    2022-09-21 12:11:27 +0000 UTC
v1.22.14    2022-09-14 22:35:04 +0000 UTC
v1.25.1     2022-09-14 19:40:59 +0000 UTC
v1.24.5     2022-09-14 16:34:14 +0000 UTC
v1.23.11    2022-09-14 16:31:19 +0000 UTC
```

2. Get details about version `v1.25.4`

```
$ ./bin/kube-score release --version v1.25.4
Kubernetes version = v1.25.4
Release time: 2022-11-10 17:02:29 +0000 UTC

********************************************************************************
kube-score release report
********************************************************************************
ImageURL                                 Tag      Digest                                                                   BuildTime                       Signed  Vulnerabilities
-----                                    -----    -----                                                                    -----                           -----   -----
registry.k8s.io/kube-apiserver           v1.25.4  sha256:ba9fc1737c5b7857f3e19183d1504ec58df0c50d970e0c008e58e8a13dc11422  2022-11-09T13:51:28.16518902Z   yes
registry.k8s.io/kube-controller-manager  v1.25.4  sha256:2526315b1c01899eab8b0fb81046083e4571d94433b293f9db124d091df98707  2022-11-09T13:51:25.456948116Z  yes
registry.k8s.io/kube-scheduler           v1.25.4  sha256:840d5b9fc29f4cddef60d832f410e3979dde2b8224cdb76dce0784394c0366a0  2022-11-09T13:51:25.462721264Z  yes
registry.k8s.io/kube-proxy               v1.25.4  sha256:1df694ba49eb1263a84c6cb32dd143d09b3e0b6cb0d48fddb3424cc4afe22e49  2022-11-09T13:51:25.801337308Z  yes
registry.k8s.io/pause                    3.8      sha256:9001185023633d17a2f98ff69b6ff2615b8ea02a825adffa40422f51dfdcde9d  2022-06-15T21:45:49.501726017Z  yes
registry.k8s.io/etcd                     3.5.5-0  sha256:b83c1d70989e1fe87583607bf5aee1ee34e52773d4755b95f5cf5a451962f3a4  2022-09-16T00:58:28.793062801Z  yes
registry.k8s.io/coredns/coredns          v1.9.3   sha256:8e352a029d304ca7431c6507b56800636c321cb52289686a581ab70aaa8a2e2a  2022-05-27T17:31:38.60727784Z   yes
```

3. Get Release details for Third-Party component.
List of third party components can be found in .kube_score.yaml file as part of the thirdPartyComponentConfig. As seen in the previous example, by default the release will pertain to Kubernetes if no components are specified. To select third party component use the flag '-component component_name'. To add support for additional thirdparty components, add the corresponding section under thirdPartyComponents and thirdPartyComponentPolicy in the same .kube_score_use.yaml.

```
./bin/kube-score release --version v3.28.4 --component calico
Release version = v3.28.4
Release time: 2025-04-15 20:45:44 +0000 UTC
Component name: calico version: v3.28.4
********************************************************************************
kube-score release report
********************************************************************************
ImageURL  Tag                                                     Digest          BuildTime                      Signed  Vulnerabilities
-----     -----                                                   -----           -----                          -----   -----
https     //github.com/projectcalico/calico/releases/tag/v3.28.4    2025-04-15 20:45:44 +0000 UTC  yes     C[1],H[1],M[5],L[1]
```

4. Get Trivy Report for specific release.
Currently support two options, use '-o /path/to/file' option to output results to a json file, or '-report true' option to write report to std out. (Reccommend using -o option)

```
./bin/kube-score release --version v3.28.4 --component calico -o /path/to/out.json
Release version = v3.28.4
Release time: 2025-04-15 20:45:44 +0000 UTC
Component name: calico version: v3.28.4
********************************************************************************
kube-score release report
********************************************************************************
ImageURL  Tag                                                     Digest          BuildTime                      Signed  Vulnerabilities
-----     -----                                                   -----           -----                          -----   -----
https     //github.com/projectcalico/calico/releases/tag/v3.28.4    2025-04-15 20:45:44 +0000 UTC  yes     C[1],H[1],M[5],L[1]
```

5. Find the upgrade recommendation for my current in-use version of `v1.23.11`. (Currently, it always uses a heuristic based approach to select best upgrade option based on vulnerabilities, Changes, EOL/EOS, etc.)

```
$ ./bin/kube-score recommend --version v1.23.11
********************************************************************************
kube-score recommendation report
********************************************************************************
Release Measures: 
        Current version: v1.23.11
        Latest version: v1.33.0
        Recommended version: v1.25.16
        Release lag (versions): 2
        Release lag (days): 1 year 8 weeks 6 days 10 hours 58 minutes 55 seconds
```

### ToDos

