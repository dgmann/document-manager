group "default" {
    targets = ["go-apps", "pdf-processor", "m1-adapter", "m1-helper", "frontend"]
}

group "cache" {
  targets = ["_go-cache"]
}

target "docker-metadata-action" {}

target "_go" {
  context = "."
  args = {
    GO_VERSION = "1.22"
  }
}

target "_go-cache" {
  inherits = ["_go"]
  dockerfile = "docker/cache.go.Dockerfile"
}

target "_go-app" {
  inherits = ["_go", "docker-metadata-action"]
  contexts = {
    cache = "target:_go-cache"
  }
}

target "go-apps" {
  name = "${service}"
  inherits = ["_go-app"]
  dockerfile = "docker/go.Dockerfile"
  matrix = {
    service = ["api", "directory-watcher", "m1-adapter"]
  }
  args = {
    SERVICE = service
  }
}

target "pdf-processor" {
  inherits = ["_go-app"]
  dockerfile = "cmd/pdf-processor/Dockerfile"
  args = {
    SERVICE = "pdf-processor"
  }
}

target "m1-helper" {
  inherits = ["_go-app"]
  dockerfile = "cmd/m1-helper/Dockerfile"
  output = ["type=local,dest=out/"]
  args = {
    SERVICE = "m1-helper"
  }
}

target "frontend" {
  inherits = ["docker-metadata-action"]
  dockerfile = "apps/frontend/Dockerfile"
  args = {
    SERVICE = "frontend"
  }
}