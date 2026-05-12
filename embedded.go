package main

import "embed"

//go:embed data/paths.json
var defaultPaths []byte

//go:embed public/*
var publicFS embed.FS
