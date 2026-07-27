// Package docs embeds the Swagger UI static files and OpenAPI YAML specs so
// they are served directly from the astrate binary without external files.
package docs

import "embed"

// SwaggerUI holds the static Swagger UI files (index.html, CSS, JS refs).
//
//go:embed swagger-ui/*
var SwaggerUI embed.FS

// API holds the OpenAPI 3.0 YAML specifications for all five API surfaces.
//
//go:embed api/*.yaml
var APIYAML embed.FS
