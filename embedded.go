package main

import "embed"

//go:embed public/index.html
var embeddedIndex []byte

//go:embed data/paths.json
var defaultPaths []byte

//go:embed public/*
var publicFS embed.FS
