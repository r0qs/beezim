package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// List of compressed websites currently maintained by Kiwix
var zims = []string{
	"gutenberg",
	"mooc",
	"other",
	"phet",
	"psiram",
	"stack_exchange",
	"ted",
	"videos",
	"vikidia",
	"wikibooks",
	"wikihow",
	"wikinews",
	"wikipedia",
	"wikiquote",
	"wikisource",
	"wikiversity",
	"wikivoyage",
	"wiktionary",
	"zimit",
}

var listWebCmd = &cobra.Command{
	Use:   "list",
	Short: "Shows a list of compressed websites currently maintained by Kiwix",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		printWebsiteList()
	},
}

func websitePath(website string) string {
	return fmt.Sprintf("%s/%s", kiwixZimURL, website)
}

func printWebsiteList() {
	const sep = "======="

	w := tabwriter.NewWriter(os.Stdout, 2, 8, 2, ' ', 0)
	fmt.Fprintf(w, "%s Kiwix Zims: %d available compressed websites %s\n", sep, len(zims), sep)
	fmt.Fprintf(w, "#\tWebsite\tURL\t\n")
	for i, site := range zims {
		fmt.Fprintf(w, "%00d\t%s\t%s\t", i+1, site, websitePath(site))
		fmt.Fprintln(w, "")
	}
	w.Flush()
}
