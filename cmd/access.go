package cmd

import (
	"log"

	"bitbucket.org/level27/lvl/types"
	"github.com/spf13/cobra"
)

// Add common commands for managing entity access to a parent command.
// entityType is the type for /{type}/{id} which this function uses.
// resolve is a function that turns an argument in the ID of the entity.
func addAccessCmds(parent *cobra.Command, entityType string, resolve func(string) int) {
	// <ENTITY> ACCESS
	var accessCmd = &cobra.Command{
		Use: "access",
		Short: "Commands for managing access to an entity",
	}

	// <ENTITY> ACCESS GET
	var accessGetCmd = &cobra.Command{
		Use: "get",
		Short: "List organisations with access to an entity",

		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			entityID := resolve(args[0])

			organisations := Level27Client.EntityGetOrganisations(entityType, entityID)

			outputFormatTableFuncs(
				organisations,
				[]string{"ID", "Name", "Type", "Members"},
				[]interface{}{"ID", "Name", "Type", func(org types.OrganisationAccess) int {
					return len(org.Users)
				}})
		},
	}

	// <ENTITY> ACCESS ADD
	var accessAddCmd = &cobra.Command{
		Use: "add",
		Short: "Grant an organisation access to an entity",

		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			entityID := resolve(args[0])
			organisationID := resolveOrganisation(args[1])

			Level27Client.EntityAddAcl(entityType,entityID, types.AclAdd{
				Organisation: organisationID,
			})

			log.Printf("Succesfully added access!")
		},
	}

	// <ENTITY> ACCESS REMOVE
	var accessRemoveCmd = &cobra.Command{
		Use: "remove",
		Short: "Revoke an organisation's access to an entity",

		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			entityID := resolve(args[0])
			organisationID := resolveOrganisation(args[1])

			Level27Client.EntityRemoveAcl(entityType, entityID, organisationID)
		
			log.Printf("%v's access removed!", args[1])
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
