package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"bitbucket.org/level27/lvl/types"
	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Commands for managing domains",
}

func init() {

	// domain command
	RootCmd.AddCommand(domainCmd)

	// Get (list of all domains)
	domainCmd.AddCommand(domainGetCmd)
	addCommonGetFlags(domainGetCmd)

	// Get details from a specific domain
	domainCmd.AddCommand(domainDescribeCmd)

	// Delete (single domain)
	domainCmd.AddCommand(domainRemoveCmd)

	// Create (single domain)
	domainCmd.AddCommand(domainCreateCmd)
	domainCreateCmd.Flags().StringVarP(&domainCreateName, "name", "n", "", "the name of the domain")
	domainCreateCmd.Flags().StringVarP(&domainCreateType, "type", "t", "", "the type of the domain")
	domainCreateCmd.Flags().StringVarP(&domainCreateLicensee, "licensee", "l", "", "the licensee of the domain")
	domainCreateCmd.Flags().StringVarP(&domainCreateOrganisation, "organisation", "o", "", "the organisation of the domain")
}

//GET LIST OF ALL DOMAINS [lvl domain get]
var domainGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a list of all current domains",
	Run: func(ccmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\t")

		domains := getDomains(args)
		for _, domain := range domains {
			fmt.Fprintln(w, strconv.Itoa(domain.ID)+"\t"+domain.Fullname+"\t"+domain.Status+"\t")
		}

		w.Flush()
	},
}

func getDomains(ids []string) []types.StructDomain {
	c := Level27Client
	if len(ids) == 0 {
		return c.Domains(optFilter, optNumber).Data
	} else {
		domains := make([]types.StructDomain, len(ids))
		for idx, id := range ids {
			domains[idx] = c.Domain("GET", id, nil).Data
		}
		return domains
	}
}

// DESCRIBE DOMAIN (get detailed info from specific domain) - [lvl domain describe <id>]
var domainDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Get detailed info about a domain",
	Run: func(cmd *cobra.Command, args []string) {
		Level27Client.DomainDescribe(args)
	},
}

// DELETE DOMAIN [lvl domain delete <id>]
var domainRemoveCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a domain",
	Run: func(cmd *cobra.Command, args []string) {
		Level27Client.DomainDelete(args)
	},
}

// CREATE DOMAIN [lvl domain create ]
var domainCreateName string
var domainCreateType string
var domainCreateLicensee string
var domainCreateOrganisation string

var domainCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new domain",

	Run: func(ccmd *cobra.Command, args []string) {
		Level27Client.DomainCreate(args)
	},
}
