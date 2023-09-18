# Installation

go install github.com/psycofdj/kustomize-yaml-helper

# Usage

```
usage: kustomize-yaml-helper [-h|--help] [-f|--file "<value>"] [-s|--stdin
                                  "<value>"] -l|--line <integer> -c|--col
                                  <integer> -a|--action
                                  (resolve|json-path|patch-path)

                                  inspect kustomization file

Arguments:

  -h  --help    Print help information
  -f  --file    path to input file
  -s  --stdin   input file name, actual content is read from stdin
  -l  --line    inspect YAML at given line
  -c  --col     inspect YAML at given column
  -a  --action  select action to perform
```

# Examples

Given `config/k8s/azure/kustomization.yaml` with this content:
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: my-namespace
resources:
- ../base
patches:
- path: patch-image-pull-secrets.yaml
  target:
    kind: Deployment
    name: ^my-.*$
- path: patch-annotation.yaml
  target:
    kind: Deployment
    name: my-deployment
- path: patch-resource.yaml
```

- get yaml value as jq path

```
prompt% kustomize-yaml-helper -f config/k8s/azure/kustomization.yaml --line 1 --col 1 --action json-path

$.apiVersion
```

- get yaml value as JSON6901 path

```
prompt% kustomize-yaml-helper -f config/k8s/azure/kustomization.yaml --line 7 --col 9 --action patch-path

/patches/0/path
```

- resolve referenced file path

```
prompt% kustomize-yaml-helper -f config/k8s/azure/kustomization.yaml --line 5 --col 7 --action resolve

<pwd>/config/k8s/base
```
