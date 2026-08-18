// Package cmd
package cmd

import (
	"fmt"

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

func install(input string) error {
	if input == "" {
		return fmt.Errorf("package name cannot be empty")
	}

	fmt.Printf("Installing %s...\n", input)

	pkg := manifest.ParseInput(input)

	manifest.AddPackage(pkg)

	return nil
}
