package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"bitbucket.org/level27/lvl/types"
	"github.com/spf13/cobra"
)

var componentCategory string
var componentType string
var componentCmd = &cobra.Command{
	Use:   "component",
	Short: "Commands related to managing app components",
}

var componentGetCmd = &cobra.Command{
	Use:   "get",

	Args: cobra.NoArgs,
	Run: func(ccmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tTYPE\tNAME\tSTATUS\t")

		components := getComponents(args)
		for _, component := range components {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", component.ID, component.AppComponentType, component.Name, component.Status)
		}

		w.Flush()
	},
}

func init() {
	RootCmd.AddCommand(componentCmd)

	componentCmd.PersistentFlags().StringVarP(&componentCategory, "category", "c", "", "Category of components to fetch")
	componentCmd.AddCommand(componentGetCmd)
	addCommonGetFlags(componentGetCmd)
	componentGetCmd.Flags().StringVarP(&componentType, "type", "t", "", "Type of components to fetch")
	componentGetCmd.MarkFlagRequired("category")
}

func getComponents(ids []string) []types.StructComponent {
	/* if len(ids) == 0 { */
	return Level27Client.Components(optFilter, optNumber, componentCategory, componentType).Components
	/* 	} else  {
		components := make([]types.StructComponent, len(ids))
		for idx, id := range ids {
			components[idx] = c.Component("GET", category, id, nil).Component
		}
		return components
	} */
}
