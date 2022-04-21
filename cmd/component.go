package cmd

// import (
// 	"bitbucket.org/level27/lvl/types"
// 	"github.com/spf13/cobra"
// )

// var componentCategory string
// var componentType string
// var componentCmd = &cobra.Command{
// 	Use:   "component",
// 	Short: "Commands related to managing app components",
// }

// var componentGetCmd = &cobra.Command{
// 	Use: "get",

// 	Args: cobra.NoArgs,
// 	Run: func(ccmd *cobra.Command, args []string) {
// 		outputFormatTable(
// 			getComponents(args),
// 			[]string{"ID", "TYPE", "NAME", "STATUS"},
// 			[]string{"ID", "AppComponentType", "Name", "Status"})
// 	},
// }

// func init() {
// 	RootCmd.AddCommand(componentCmd)

// 	componentCmd.PersistentFlags().StringVarP(&componentCategory, "category", "c", "", "Category of components to fetch")
// 	componentCmd.AddCommand(componentGetCmd)
// 	addCommonGetFlags(componentGetCmd)
// 	componentGetCmd.Flags().StringVarP(&componentType, "type", "t", "", "Type of components to fetch")
// 	componentGetCmd.MarkFlagRequired("category")
// }

// func getComponents(ids []string) []types.AppComponent2 {
// 	return Level27Client.Components(componentCategory, componentType, optGetParameters)
// }

