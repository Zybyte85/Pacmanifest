package manifest

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Zybyte85/mypm/internal/config"
)

type Package struct {
	Name    string
	Version string
	IsAUR   bool
}

func AddPackage(pkg Package) error {
	if err := os.MkdirAll(config.GetConfigPath(), 0o755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}
	f, err := os.OpenFile(filepath.Join(config.GetConfigPath(), "manifest.pkgs"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("could not open manifest: %w", err)
	}
	defer f.Close()

	var pkgStr string
	if pkg.IsAUR {
		pkgStr = "aur:" + pkg.Name
	} else {
		pkgStr = pkg.Name
	}
	if pkg.Version != "latest" {
		pkgStr += "=" + pkg.Version
	}
	pkgStr += "\n"

	if _, err := f.WriteString(pkgStr); err != nil {
		return fmt.Errorf("could not write to manifest: %w", err)
	}

	return nil
}

func ParseInput(input string) Package {
	pkg := Package{}

	if strings.HasPrefix(input, "aur:") {
		pkg.IsAUR = true
		input = strings.TrimPrefix(input, "aur:")
	}

	var found bool
	pkg.Name, pkg.Version, found = strings.Cut(input, "=")

	if !found {
		pkg.Version = "latest"
	}

	return pkg
}

func ParseManifest(manifestPath string) ([]Package, error) {
	pkgs := []Package{}

	f, err := os.Open(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("could not open manifest: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		pkgs = append(pkgs, ParseInput(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("could not read manifest: %w", err)
	}

	fmt.Println(pkgs)

	return nil, nil
}
