// Package templates provides embedded HTML templates for the application.
package templates

import "embed"

//go:embed *.html
var FS embed.FS
