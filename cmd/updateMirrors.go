package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// updateMirrorsCmd represents the updateMirrors command
var updateMirrorsCmd = &cobra.Command{
	Use:   "updateMirrors",
	Short: "Updaqte the mirrorlist to check for newer repos.",

	RunE: func(cmd *cobra.Command, args []string) error {
		err := getBestMirrors(country)
		if err != nil {
			return fmt.Errorf("error getting mirrors: %w", err)
		}

		return nil
	},
}

var country string

func init() {
	rootCmd.AddCommand(updateMirrorsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateMirrorsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateMirrorsCmd.Flags().String("country", &country, "")
	updateMirrorsCmd.Flags().StringVarP(&country, "country", "c", "", "Set the country to be used for fetching mirrors (use two letter code)")
}

type Mirrors struct {
	URLs []URL `json:"urls"`
}

type URL struct {
	URL         string  `json:"url"`
	Protocol    string  `json:"protocol"`
	Completion  float64 `json:"completion_pct"`
	Delay       int     `json:"delay"`
	Score       float64 `json:"score"`
	Active      bool    `json:"active"`
	CountryCode string  `json:"country_code"`
}

func FilterMirrors(mirrors []URL, targetCountry string) []URL {
	var filtered []URL
	for _, m := range mirrors {
		// If the mirror is active, fully complete, uses HTTPS, is in the user's country, and current.
		if m.Active &&
			m.Completion == 1.0 &&
			m.Protocol == "https" &&
			m.CountryCode == targetCountry &&
			m.Delay < 21600 {
			filtered = append(filtered, m)
		}
	}
	return filtered
}

func getBestMirrors(targetCountry string) error {
	resp, err := http.Get("https://archlinux.org/mirrors/status/json/")
	if err != nil {
		return fmt.Errorf("could not fetch mirrors")
	}
	defer resp.Body.Close()

	var data Mirrors
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	if targetCountry == "" {
		targetCountry = getRegion()
	}

	results := FilterMirrors(data.URLs, targetCountry)

	// Sort mirrors by their "score" value
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score < results[j].Score
	})

	// TODO: Eventually, we want to test the speed of each mirror and rank them in order

	// Check if we actually found anything
	if len(results) == 0 {
		return fmt.Errorf("no countries found for country: %s", targetCountry)
	}

	var sb strings.Builder
	for _, m := range results {
		sb.WriteString(m.URL + "\n")
	}
	return os.WriteFile("mirrorlist", []byte(sb.String()), 0o644)
}

func getRegion() string {
	// Service for getting the current country
	url := "https://ipapi.co/country/"

	resp, err := http.Get(url)
	if err != nil {
		// Fallback to US
		return "US"
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	res := strings.TrimSpace(string(body))

	if len(res) == 2 {
		return res
	}
	return "US"
}
