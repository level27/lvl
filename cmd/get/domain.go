package get

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"bitbucket.org/level27/lvl/cmd"
	"bitbucket.org/level27/lvl/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var domainGetCmd = &cobra.Command{
	Use:   "domain [IDs to retrieve]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ArbitraryArgs,
	Run: func(ccmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\t")

		domains := getDomains(args)
		for _, domain := range domains {
			fmt.Fprintln(w, strconv.Itoa(domain.ID)+"\t"+domain.Fullname+"\t"+domain.Status+"\t")
		}
	
		w.Flush()
	},
}

func init() {
	GetCmd.AddCommand(domainGetCmd)
}

func getDomains(ids []string) []types.StructDomain {
	c := cmd.Level27Client
	if len(ids) == 0 {
		numberToGet := viper.GetString("number")
		domainFilter := viper.GetString("filter")
	
		return c.Domains(domainFilter, numberToGet).Data
	} else {
		domains := make([]types.StructDomain, len(ids))
		for idx, id := range ids {
			domains[idx] = c.Domain("GET", id, nil).Data
		}
		return domains
	}
}
