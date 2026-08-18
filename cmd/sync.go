/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
		installFromManifest()
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

func installFromManifest() {
	var pacmanPkgs []string
	var aurPkgs []string

	manifestPath := filepath.Join(config.GetConfigPath(), "manifest.pkgs")
	pkgs, err := manifest.ParseManifest(manifestPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, pkg := range pkgs {
		if pkg.IsAUR {
			// If the package is an AUR package, use the AUR helper to install
			aurPkgs = append(aurPkgs, pkg.Name)

		} else if pkg.Version == "latest" {
			// Not AUR, and latest version, so install normally
			pacmanPkgs = append(pacmanPkgs, pkg.Name)
		}
	}
	commandInstall(pacmanPkgs, "pacman", true)
	commandInstall(aurPkgs, aurHelper, false)
}

func commandInstall(pkgs []string, command string, sudo bool) {
	if len(pkgs) == 0 {
		// Nothing to install
		return
	}

	// Build the argument list for the package manager
	args := []string{"-S", "--noconfirm", "--needed"}
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
