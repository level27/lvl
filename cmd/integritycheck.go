package cmd

import (
	"fmt"

	"bitbucket.org/level27/lvl/types"
	"bitbucket.org/level27/lvl/utils"
	"github.com/spf13/cobra"
)

func outputFormatIntegrityCheckTable(checks []types.IntegrityCheck) {
	outputFormatTableFuncs(
		checks,
		[]string{"ID", "STATUS", "DATE"},
		[]interface{}{"Id", "Status", func(s types.IntegrityCheck) string {
			return utils.FormatUnixTime(s.DtRequested)
		}})
}

// Add common commands for managing integrity checks to a parent command.
// entityType is the type for /{type}/{id} which this function uses.
// resolve is a function that turns an argument in the ID of the entity.
func addIntegrityCheckCmds(parent *cobra.Command, entityType string, resolve func(string) int) {
	var integrityCmd = &cobra.Command{
		Use:   "integrity",
		Short: "Commands for managing integrity checks",
	}

	var integrityGetCmd = &cobra.Command{
		Use:   "get [entity]",
		Short: "Get a list of all integrity checks for an entity",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			entityID := resolve(args[0])

			checks := resolveGets(
				// First arg is entity ID.
				args[1:],
				// Can't do lookups for integrity checks.
				func(name string) []types.IntegrityCheck { return nil },
				// Large funcs to pass entity type and ID along.
				func(checkID int) types.IntegrityCheck { return Level27Client.EntityIntegrityCheck(entityType, entityID, checkID)},
				func(get types.CommonGetParams) []types.IntegrityCheck { return Level27Client.EntityIntegrityChecks(entityType, entityID, get)})

			outputFormatIntegrityCheckTable(checks)
		},
	}

	var integrityCheckDoJobs bool = true
	var integrityCheckForceJobs bool = false
	var integrityCreateCmd = &cobra.Command{
		Use:   "create [entity]",
		Short: "Create a new integrity report",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			entityID := resolve(args[0])

			result := Level27Client.EntityIntegrityCreate(entityType, entityID, integrityCheckDoJobs, integrityCheckForceJobs)
			outputFormatTemplate(result, "templates/integrityCreate.tmpl")
		},
	}

	var integrityDownload string
	var integrityDownloadCmd = &cobra.Command{
		Use:   "download [entity] [check id]",
		Short: "Download an integrity check as PDF file",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			entityID := resolve(args[0])
			checkId, err := convertStringToId(args[1])
			cobra.CheckErr(err)

			if integrityDownload == "" {
				// Auto-generate file name.
				integrityDownload = fmt.Sprintf("integritycheck_%d_%s_%d.pdf", checkId, entityType, entityID)
			}

			Level27Client.EntityIntegrityCheckDownload(entityType, entityID, checkId, integrityDownload)
		},
	}


	parent.AddCommand(integrityCmd)

	integrityCmd.AddCommand(integrityGetCmd)
	addCommonGetFlags(integrityGetCmd)

	integrityCmd.AddCommand(integrityCreateCmd)
	flags := integrityCreateCmd.Flags()
	flags.BoolVar(&integrityCheckDoJobs, "doJobs", integrityCheckDoJobs, "Create jobs")
	flags.BoolVar(&integrityCheckForceJobs, "forceJobs", integrityCheckForceJobs, "Create jobs even if integrity check failed")

	integrityCmd.AddCommand(integrityDownloadCmd)
	integrityDownloadCmd.Flags().StringVarP(&integrityDownload, "file", "f", "", "File to download the report to. This defaults to a generated file name in the current directory.")
}
