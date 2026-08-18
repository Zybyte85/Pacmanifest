/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/Zybyte85/mypm/internal/config"
	"github.com/Zybyte85/mypm/internal/manifest"
	"github.com/spf13/cobra"
)

var aurHelper string

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Reads the manifest and applies changes",
	Run: func(cmd *cobra.Command, args []string) {
		err := syncRepos()
		if err != nil {
			fmt.Println(err)
			return
		}

		err = installFromManifest()
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	syncCmd.Flags().StringVarP(&aurHelper, "aur-helper", "a", "yay", "AUR helper to use")
}

func installFromManifest() error {
	var pacmanPkgs []string
	var aurPkgs []string
	var archivePkgs []manifest.Package

	manifestPath := filepath.Join(config.GetConfigPath(), "manifest.pkgs")
	pkgs, err := manifest.ParseManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("could not parse manifest: %v", err)
	}
	for _, pkg := range pkgs {
		var latestVersion string

		if pkg.Version != "latest" {
			latestVersion, err = seeLatestVersionOfPackage(pkg.Name)

			if err != nil {
				return fmt.Errorf("could not see latest version of package %s: %v", pkg.Name, err)
			}
		}

		if pkg.IsAUR {
			// If the package is an AUR package, use the AUR helper to install
			aurPkgs = append(aurPkgs, pkg.Name)

		} else if pkg.Version == "latest" || pkg.Version == latestVersion {
			// Not AUR, and latest version, so install normally
			pacmanPkgs = append(pacmanPkgs, pkg.Name)
		} else {
			// Not AUR, and specific version, so install from archive
			archivePkgs = append(archivePkgs, pkg)
		}
	}
	runInstallCommand(pacmanPkgs, "pacman", []string{"-S", "--noconfirm", "--needed"}, true)
	runInstallCommand(aurPkgs, aurHelper, []string{"-S", "--noconfirm", "--needed"}, false)
	archiveInstall(archivePkgs)

	return nil
}

func seeLatestVersionOfPackage(pkgName string) (string, error) {
	cmd := exec.Command("pacman", "-Si", pkgName)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Failed to get latest package version: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	_, versionStr, _ := strings.Cut(lines[2], ": ")
	// Split the version string to remove the release number
	latestVersion, _, _ := strings.Cut(versionStr, "-")

	return string(latestVersion), nil
}

func runInstallCommand(pkgs []string, command string, args []string, sudo bool) {
	if len(pkgs) == 0 {
		// Nothing to install
		return
	}

	// Build the argument list for the package manager
	args = append(args, pkgs...)

	var cmd *exec.Cmd
	if sudo {
		cmd = exec.Command("sudo", append([]string{command}, args...)...)
	} else {
		cmd = exec.Command(command, args...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Show what will be installed
	fmt.Printf("Installing %s with %s\n", strings.Join(pkgs, ", "), command)
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}

func archiveInstall(pkgs []manifest.Package) error {
	if len(pkgs) == 0 {
		// Nothing to install
		return fmt.Errorf("No packages to install")
	}

	var architectureStr string
	switch runtime.GOARCH {
	case "amd64":
		architectureStr = "x86_64"
	case "arm64":
		architectureStr = "aarch64"
	default:
		architectureStr = "x86_64" // Fallback to x86_64
	}

	var filePaths []string
	cachePath := filepath.Join(config.GetConfigPath(), "cache")

	for _, pkg := range pkgs {
		// Fetch archive list to find package release number
		fmt.Printf("Downloading %s from archive\n", pkg.Name)
		baseUrl := fmt.Sprintf("https://archive.archlinux.org/packages/%s/%s", pkg.Name[:1], pkg.Name)

		resp, err := http.Get(baseUrl)
		if err != nil {
			return fmt.Errorf("Error fetching directory for %s: %v", pkg.Name, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Package %s not found on archive (Status: %s)", pkg.Name, resp.Status)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error reading response body for %s: %v", pkg.Name, err)
		}
		htmlContent := string(bodyBytes)

		pattern := fmt.Sprintf(`%s-%s-([a-zA-Z0-9.]+)-(?:%s|any)\.pkg\.tar\.zst`,
			regexp.QuoteMeta(pkg.Name),
			regexp.QuoteMeta(pkg.Version),
			architectureStr,
		)
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(htmlContent)

		if len(matches) < 2 {
			fmt.Printf("Could not find a matching archive file for %s version %s\n", pkg.Name, pkg.Version)
			continue
		}

		// pkgrel := matches[1]
		exactFileName := matches[0]

		filePath := filepath.Join(cachePath, exactFileName)
		filePaths = append(filePaths, filePath)

		downloadUrl := fmt.Sprintf("%s/%s", baseUrl, exactFileName)

		err = os.MkdirAll(cachePath, 0o755)
		if err != nil {
			return fmt.Errorf("Error creating cache directory: %v", err)
		}

		downloadFile(downloadUrl, filePath)
	}

	runInstallCommand(filePaths, "pacman", []string{"-U", "--noconfirm", "--needed"}, true)

	for _, filePath := range filePaths {
		os.Remove(filePath)
	}

	return nil
}
func downloadFile(url, dest string) error {
	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("Error creating file: %v", err)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error downloading file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to download file (Status: %d)", resp.StatusCode)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("Error saving file: %v", err)
	}

	return nil
}

func syncRepos() error {
	cmd := exec.Command("sudo", "pacman", "-Sy")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not sync repositories: %v", err)
	}

	return nil
}
