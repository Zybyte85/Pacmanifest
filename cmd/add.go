// Package cmd
package cmd

import (
	"fmt"

	"github.com/Zybyte85/mypm/internal/manifest"
	"github.com/spf13/cobra"
)

// addCmd represents the install command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds the specified package to the manifest",
	Args:  cobra.MinimumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		return install(args[0])
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func install(input string) error {
	pkg := manifest.ParseInput(input)

	err := manifest.AddPackage(pkg)
	if err != nil {
		return fmt.Errorf("could not add '%s' to manifest: %v", pkg.Name, err)
	}

	fmt.Printf("Added '%s' to manifest. Run 'pacmanifest sync' to install.\n", pkg.Name)

	return nil
}
