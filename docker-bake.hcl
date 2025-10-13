group "default" {
  targets = ["go-backend", "documents-service"]
}

target "go-backend" {
  context = "./"
  dockerfile = "Dockerfile"
  args = {
    GOPROXY = "https://proxy.golang.org,direct"
  }
}

target "document-service" {
  context = "./services"
}