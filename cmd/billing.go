package cmd

import (
	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

// Add common commands for managing entity billing to a parent command.
// entityType is the type for /{type}/{id} which this function uses.
// resolve is a function that turns an argument in the ID of the entity.
func addBillingCmds(parent *cobra.Command, entityType string, resolve func(string) (l27.IntID, error)) {
	// BILLING
	billingCmd := &cobra.Command{
		Use:   "billing",
		Short: "Manage entity's invoicing (BillableItem)",
	}

	var billingOnExternalInfo string

	// BILLING ON
	billingOnCmd := &cobra.Command{
		Use:   "on [domain] [flags]",
		Short: "Turn on billing for an entity (admin only)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entityID, err := resolve(args[0])
			if err != nil {
				return err
			}

			req := l27.BillPostRequest{
				ExternalInfo: billingOnExternalInfo,
			}

			err = Level27Client.EntityBillableItemCreate(entityType, entityID, req)
			if err != nil {
				return err
			}

			return nil
		},
	}

	billingOnCmd.Flags().StringVarP(&billingOnExternalInfo, "externalinfo", "e", "", "ExternalInfo (required when billableitemInfo entities for an Organisation exist in db)")

	// BILLING OFF
	billingOffCmd := &cobra.Command{
		Use:   "off [domainID]",
		Short: "Turn off the billing for an entity (admin only)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entityID, err := resolve(args[0])
			if err != nil {
				return err
			}

			Level27Client.EntityBillableItemDelete(entityType, entityID)
			return nil
		},
	}

	billingCmd.AddCommand(billingOffCmd)
	billingCmd.AddCommand(billingOnCmd)

	parent.AddCommand(billingCmd)
}
