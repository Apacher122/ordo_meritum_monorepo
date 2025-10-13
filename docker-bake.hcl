group "default" {
  targets = ["go-api", "document-service"]
}

target "go-api" {
  context = "./ordo_meritum_go_server"
  dockerfile = "Dockerfile"
  args = {
    GOPROXY = "https://proxy.golang.org,direct"
  }
}

target "document-service" {
  context = "./document-service"
}