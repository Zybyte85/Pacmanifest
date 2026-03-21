// Package cmd
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Zybyte85/mypm/internal/manifest"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs the specified package",
	Args:  cobra.MinimumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		err := install(args[0])
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func install(pkg string) error {
	if pkg == "" {
		return fmt.Errorf("package name cannot be empty")
	}

	fmt.Printf("Installing %s...\n", pkg)

	fmt.Println(parseInput(pkg))

	f, err := os.OpenFile(filepath.Join(cfgPath, "manifest.pkgs"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("could not open manifest: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(pkg + "\n"); err != nil {
		return fmt.Errorf("could not write to manifest: %w", err)
	}

	return nil
}

func parseInput(input string) (Package manifest.Package) {
	fmt.Errorf("hhi")
	pkg := manifest.Package{}

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
