package cmd

import (
	"errors"
	"fmt"

	"github.com/level27/l27-go"
	"github.com/level27/lvl/utils"
	"github.com/spf13/cobra"
)

func outputFormatIntegrityCheckTable(checks []l27.IntegrityCheck) {
	outputFormatTableFuncs(
		checks,
		[]string{"ID", "STATUS", "DATE"},
		[]interface{}{"ID", "Status", func(s l27.IntegrityCheck) string {
			return utils.FormatUnixTime(s.DtRequested)
		}})
}

// Add common commands for managing integrity checks to a parent command.
// entityType is the type for /{type}/{id} which this function uses.
// resolve is a function that turns an argument in the ID of the entity.
func addIntegrityCheckCmds(parent *cobra.Command, entityType string, resolve func(string) (l27.IntID, error)) {
	var integrityCmd = &cobra.Command{
		Use:   "integrity",
		Short: "Commands for managing integrity checks",
	}

	var integrityGetCmd = &cobra.Command{
		Use:   "get [entity]",
		Short: "Get a list of all integrity checks for an entity",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entityID, err := resolve(args[0])
			if err != nil {
				return err
			}

			checks, err := resolveGets(
				// First arg is entity ID.
				args[1:],
				// Can't do lookups for integrity checks.
				func(name string) ([]l27.IntegrityCheck, error) {
					return nil, errors.New("integrity checks cannot be looked up by any name")
				},
				// Large funcs to pass entity type and ID along.
				func(checkID l27.IntID) (l27.IntegrityCheck, error) {
					return Level27Client.EntityIntegrityCheck(entityType, entityID, checkID)
				},
				func(get l27.CommonGetParams) ([]l27.IntegrityCheck, error) {
					return Level27Client.EntityIntegrityChecks(entityType, entityID, get)
				})

			if err != nil {
				return err
			}

			outputFormatIntegrityCheckTable(checks)
			return nil
		},
	}

	var integrityCheckDoJobs bool = true
	var integrityCheckForceJobs bool = false
	var integrityCreateCmd = &cobra.Command{
		Use:   "create [entity]",
		Short: "Create a new integrity report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entityID, err := resolve(args[0])
			if err != nil {
				return err
			}

			result, err := Level27Client.EntityIntegrityCreate(entityType, entityID, integrityCheckDoJobs, integrityCheckForceJobs)
			if err != nil {
				return err
			}

			if optWait {
				result, err = waitForStatus(
					func() (l27.IntegrityCheck, error) {
						return Level27Client.EntityIntegrityCheck(entityType, entityID, result.ID)
					},
					func(s l27.IntegrityCheck) string { return s.Status },
					"ok",
					[]string{"to_create", "creating"},
				)

				if err != nil {
					return fmt.Errorf("waiting on check status failed: %s", err.Error())
				}
			}

			outputFormatTemplate(result, "templates/integrityCreate.tmpl")
			return nil
		},
	}

	var integrityDownload string
	var integrityDownloadCmd = &cobra.Command{
		Use:   "download [entity] [check id]",
		Short: "Download an integrity check as PDF file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			entityID, err := resolve(args[0])
			if err != nil {
				return err
			}

			checkID, err := convertStringToID(args[1])
			if err != nil {
				return err
			}

			err = Level27Client.EntityIntegrityCheckDownload(entityType, entityID, checkID, integrityDownload)
			if err != nil {
				return err
			}

			return nil
		},
	}

	parent.AddCommand(integrityCmd)

	integrityCmd.AddCommand(integrityGetCmd)
	addCommonGetFlags(integrityGetCmd)

	integrityCmd.AddCommand(integrityCreateCmd)
	addWaitFlag(integrityCreateCmd)
	flags := integrityCreateCmd.Flags()
	flags.BoolVar(&integrityCheckDoJobs, "doJobs", integrityCheckDoJobs, "Create jobs")
	flags.BoolVar(&integrityCheckForceJobs, "forceJobs", integrityCheckForceJobs, "Create jobs even if integrity check failed")

	integrityCmd.AddCommand(integrityDownloadCmd)
	integrityDownloadCmd.Flags().StringVarP(&integrityDownload, "file", "f", "", "File to download the report to. This defaults to a generated file name in the current directory.")
}
