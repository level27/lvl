package cmd

import (
	"fmt"
	"log"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

// Add common commands for managing entity access to a parent command.
// entityType is the type for /{type}/{id} which this function uses.
// resolve is a function that turns an argument in the ID of the entity.
func addAccessCmds(parent *cobra.Command, entityType string, resolve func(string) (int, error)) {
	// <ENTITY> ACCESS
	var accessCmd = &cobra.Command{
		Use:   "access",
		Short: "Commands for managing access to an entity",
	}

	// <ENTITY> ACCESS GET
	var accessGetCmd = &cobra.Command{
		Use:   "get",
		Short: "List organisations with access to an entity",

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entityID, err := resolve(args[0])
			if err != nil {
				return fmt.Errorf("unable to resolve %s '%s': %v", entityType, args[0], err)
			}

			organisations, err := Level27Client.EntityGetOrganisations(entityType, entityID)
			if err != nil {
				return err
			}

			outputFormatTableFuncs(
				organisations,
				[]string{"ID", "Name", "Type", "Members"},
				[]interface{}{"ID", "Name", "Type", func(org l27.OrganisationAccess) int {
					return len(org.Users)
				}})

			return nil
		},
	}

	// <ENTITY> ACCESS ADD
	var accessAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Grant an organisation access to an entity",

		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			entityID, err := resolve(args[0])
			if err != nil {
				return err
			}

			organisationID, err := resolveOrganisation(args[1])
			if err != nil {
				return err
			}

			_, err = Level27Client.EntityAddAcl(entityType, entityID, l27.AclAdd{
				Organisation: organisationID,
			})

			if err == nil {
				log.Printf("Succesfully added access!")
			}

			return err
		},
	}

	// <ENTITY> ACCESS REMOVE
	var accessRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Revoke an organisation's access to an entity",

		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			entityID, err := resolve(args[0])
			if err != nil {
				return err
			}

			organisationID, err := resolveOrganisation(args[1])
			if err != nil {
				return err
			}

			err = Level27Client.EntityRemoveAcl(entityType, entityID, organisationID)

			if err == nil {
				log.Printf("%v's access removed!", args[1])
			}

			return err
		},
	}

	parent.AddCommand(accessCmd)

	// <ENTITY> ACCESS GET
	accessCmd.AddCommand(accessGetCmd)

	// <ENTITY> ACCESS ADD
	accessCmd.AddCommand(accessAddCmd)

	// <ENTITY> ACCESS REMOVE
	accessCmd.AddCommand(accessRemoveCmd)

}
