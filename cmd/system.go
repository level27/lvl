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
	flags.StringVarP(&systemCreateFqdn, "Fqdn", "", "", "Valid hostname for the system")
	flags.StringVarP(&systemCreateRemarks, "remarks", "", "", "Remarks (Admin only)")
	flags.IntVarP(&systemCreateDisk, "disk", "", 0, "Disk (non-editable)")
	flags.IntVarP(&systemCreateDisk, "cpu", "", 0, "Cpu (Required for Level27 systems)")
	flags.IntVarP(&systemCreateDisk, "memory", "", 0, "Memory (Required for Level27 systems)")
	flags.StringVarP(&systemCreateManageType, "management", "", "basic", "Managament type (default: basic)")
	flags.BoolVarP(&systemCreatePublicNetworking, "publicNetworking", "", true, "For digitalOcean servers always true. (non-editable)")
	flags.IntVarP(&systemCreateImage, "image", "", 0, "The ID of a systemimage. (must match selected configuration and zone. non-editable)")
	flags.IntVarP(&systemCreateOrganisation, "organisation", "", 0, "The unique ID of an organisation")
	flags.IntVarP(&systemCreateProviderConfig, "provider", "", 0, "The unique ID of a SystemproviderConfiguration")
	flags.IntVarP(&systemCreateZone, "zone", "", 0, "The unique ID of a zone")
	flags.IntVarP(&systemCreateSecurityUpdates, "security", "", 0, "installSecurityUpdates (default: random POST:1-8, PUT:0-12)")
	flags.StringVarP(&systemCreateAutoTeams, "autoTeams", "", "", "A csv list of team ID's")
	flags.StringVarP(&systemCreateExternalInfo, "externalInfo", "", "", "ExternalInfo (required when billableItemInfo entities for an organisation exist in db)")
	flags.IntVarP(&systemCreateOperatingSystemVersion, "version", "", 0, "The unique ID of an OperatingsystemVersion (non-editable)")
	flags.IntVarP(&systemCreateParentSystem, "parent", "", 0, "The unique ID of a system (parent system)")
	flags.StringVarP(&systemCreateType, "type", "", "", "System type")
	flags.StringArrayP("networks", "", []string{""}, "Array of network IP's. (default: null)")

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
var systemCreateDisk, systemCreateCpu, systemCreateMemory int
var systemCreateManageType string
var systemCreatePublicNetworking bool
var systemCreateImage, systemCreateOrganisation, systemCreateProviderConfig, systemCreateZone int
var systemCreateSecurityUpdates int
var systemCreateAutoTeams, systemCreateExternalInfo string
var systemCreateOperatingSystemVersion, systemCreateParentSystem int
var systemCreateType string
var systemCreateAutoNetworks []interface{}
var managementTypeArray = []string{"basic", "professional", "enterprise", "professional_level27"}
var systemCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new system",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
