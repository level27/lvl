package get

import (
	"fmt"
	"os"
	"text/tabwriter"

	"bitbucket.org/level27/lvl/cmd"
	"bitbucket.org/level27/lvl/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var componentCategory string
var componentType string
var componentGetCmd = &cobra.Command{
	Use:   "component",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	componentGetCmd.Flags().StringVarP(&componentCategory, "category", "c", "", "Category of components to fetch")
	componentGetCmd.Flags().StringVarP(&componentType, "type", "t", "", "Type of components to fetch")
	componentGetCmd.MarkFlagRequired("category")
	viper.BindPFlag("category", componentGetCmd.Flags().Lookup("category"))
	viper.BindPFlag("type", componentGetCmd.Flags().Lookup("type"))

	GetCmd.AddCommand(componentGetCmd)
}

func getComponents(ids []string) []types.StructComponent {
	c := cmd.Level27Client
	category := viper.GetString("category")
	/* if len(ids) == 0 { */
		numberToGet := viper.GetString("number")
		filter := viper.GetString("filter")
		cType := viper.GetString("type")
		return c.Components(filter, numberToGet, category, cType).Components
/* 	} else  {
		components := make([]types.StructComponent, len(ids))
		for idx, id := range ids {
			components[idx] = c.Component("GET", category, id, nil).Component
		}
		return components
	} */
}
