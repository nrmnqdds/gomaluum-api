package main

import (
	"embed"

	cmd "github.com/nrmnqdds/gomaluum-api/cmd"
	"github.com/nrmnqdds/gomaluum-api/cmd/application"
)

//go:embed docs/*
var SwaggerDocsPath embed.FS

func main() {
	application.SwaggerDocsPath = &SwaggerDocsPath
	cmd.StartApplication()
}
