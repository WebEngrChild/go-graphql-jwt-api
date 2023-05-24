# Getting Started

```shell:
$ make init
$ make migrate
$ make start
```

## Build for ECR

```
# project root
$ docker build --no-cache --target runner -t api-dev-repo --platform linux/amd64 -f ./build/docker/go/Dockerfile .
```

## Library & Tool 
### Dockle

```
$ brew install goodwithtech/r/dockle
$ dockle api-dev-repo
```

### trivy

```
$ brew install aquasecurity/trivy/trivy
$ trivy image --severity CRITICAL --ignore-unfixed api-dev-repo
```

### hadolint

```
$ brew install hadolint
$ hadolint --ignore DL3018 build/docker/go/Dockerfile
```