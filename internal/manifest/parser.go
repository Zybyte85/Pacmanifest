// Package manifest provides functions for working with the manifest files
package manifest

type Package struct {
	Name    string
	Version string
	IsAUR   bool
}
