package cmd

import (
	"log"

	"bitbucket.org/level27/lvl/types"
	"github.com/spf13/cobra"
)

// MAIN COMMAND
var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Commands for managing systems",
}

func init() {
	//adding main command to root
	RootCmd.AddCommand(systemCmd)

	//Toplevel subcommands (get/post)
	// --- GET
	systemCmd.AddCommand(systemGetCmd)
	addCommonGetFlags(systemGetCmd)

	// --- CREATE
	systemCmd.AddCommand(systemCreateCmd)
	flags := systemCreateCmd.Flags()
	flags.StringVarP(&systemCreateName, "name", "n", "", "The name you want to give the system")
	flags.StringVarP(&systemCreateFqdn, "Fqdn" ,"","", "Valid hostname for the system")
	flags.StringVarP(&systemCreateRemarks, "remarks", "", "", "Remarks (Admin only)")
}


//------------------------------------------------- SYSTEM TOPLEVEL (GET / CREATE) ----------------------------------
//----------------------------------------- GET ---------------------------------------
var systemGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get a list of all curent systems",
	Run: func(cmd *cobra.Command, args []string) {
		ids, err := convertStringsToIds(args)
		if err != nil {
			log.Fatalln("Invalid domain ID")
		}
		outputFormatTable(getSystems(ids), []string{"ID", "NAME", "STATUS"}, []string{"Id", "Name", "Status"})

	},
}

func getSystems(ids []int) []types.System {

	if len(ids) == 0 {
		return Level27Client.SystemGetList(optGetParameters)
	} else {
		systems := make([]types.System, len(ids))
		for idx, id := range ids {
			systems[idx] = Level27Client.SystemGetSingle(id)
		}
		return systems
	}

}


//----------------------------------------- CREATE ---------------------------------------
// vars needed to save flag data.
var systemCreateName, systemCreateFqdn, systemCreateRemarks string 

var systemCreateCmd = &cobra.Command{
	Use: "create",
	Short: "Create a new system",
	Run: func(cmd *cobra.Command, args []string) {


	},
}