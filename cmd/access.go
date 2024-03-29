package cmd

import (
	"fmt"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

// Add common commands for managing entity access to a parent command.
// entityType is the type for /{type}/{id} which this function uses.
// resolve is a function that turns an argument in the ID of the entity.
func addAccessCmds(parent *cobra.Command, entityType string, resolve func(string) (l27.IntID, error)) {
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

			acl, err := Level27Client.EntityAddAcl(entityType, entityID, l27.AclAdd{
				Organisation: organisationID,
			})

			if err != nil {
				return err
			}

			outputFormatTemplate(acl, "templates/entities/acl/added.tmpl")

			return nil
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
			if err != nil {
				return err
			}

			outputFormatTemplate(nil, "templates/entities/acl/removed.tmpl")

			return nil
		},
	}

	parent.AddCommand(accessCmd)

	// <ENTITY> ACCESS GET
	accessCmd.AddCommand(accessGetCmd)

	// <ENTITY> ACCESS ADD
	accessCmd.AddCommand(accessAddCmd)

	// <ENTITY> ACCESS REMOVE
	accessCmd.AddCommand(accessRemoveCmd)

	// <ENTITY> TEAM
	var teamCmd = &cobra.Command{
		Use:   "team",
		Short: "Commands for managing teams that have access to an entity",
	}

	// <ENTITY> TEAM ADD
	var teamAddCmd = &cobra.Command{
		Use:   "add <entity> <organisation> <team>",
		Short: "Add a team to this entity",

		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			entityID, err := resolve(args[0])
			if err != nil {
				return err
			}

			organisationID, err := resolveOrganisation(args[1])
			if err != nil {
				return err
			}

			teamID, err := resolveOrganisationTeam(organisationID, args[2])
			if err != nil {
				return err
			}

			teamEntity, err := Level27Client.OrganisationTeamEntityAdd(organisationID, teamID, entityType, entityID)
			if err != nil {
				return err
			}

			outputFormatTemplate(teamEntity, "templates/entities/organisationTeamEntity/add.tmpl")

			return nil
		},
	}

	// <ENTITY> TEAM REMOVE
	var teamRemoveCmd = &cobra.Command{
		Use:   "remove <entity> <organisation> <team>",
		Short: "Add a team to this entity",

		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			entityID, err := resolve(args[0])
			if err != nil {
				return err
			}

			organisationID, err := resolveOrganisation(args[1])
			if err != nil {
				return err
			}

			teamID, err := resolveOrganisationTeam(organisationID, args[2])
			if err != nil {
				return err
			}

			err = Level27Client.OrganisationTeamEntityRemove(organisationID, teamID, entityType, entityID)
			if err != nil {
				return err
			}

			outputFormatTemplate(nil, "templates/entities/organisationTeamEntity/remove.tmpl")

			return nil
		},
	}

	parent.AddCommand(teamCmd)

	teamCmd.AddCommand(teamAddCmd)

	teamCmd.AddCommand(teamRemoveCmd)
}
