package cmd

import (
	"fmt"
	"log"
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

	// ----------------- DOMAINS ------------------------
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
	//required flags
	domainCreateCmd.Flags().StringVarP(&domainCreateName, "name", "n", "", "the name of the domain (Required)")	
	domainCreateCmd.Flags().IntVarP(&domainCreateType, "type", "t", 0, "the type of the domain (Required)")
	domainCreateCmd.Flags().IntVarP(&domainCreateLicensee, "licensee", "l", 0, "The unique identifier of a domaincontact with type licensee (Required)")
	domainCreateCmd.Flags().IntVarP(&domainCreateOrganisation, "organisation", "o", 0, "the organisation of the domain (Required)")
	domainCreateCmd.MarkFlagRequired("name")
	domainCreateCmd.MarkFlagRequired("type")
	domainCreateCmd.MarkFlagRequired("licensee")
	domainCreateCmd.MarkFlagRequired("organisation")

	// optional flags
	domainCreateCmd.Flags().StringVarP(&domainCreateNs1, "nameserver1", "", "", "Nameserver")
	domainCreateCmd.Flags().StringVarP(&domainCreateNs2, "nameserver2", "", "", "Nameserver")
	domainCreateCmd.Flags().StringVarP(&domainCreateNs3, "nameserver3", "", "", "Nameserver")
	domainCreateCmd.Flags().StringVarP(&domainCreateNs4, "nameserver4", "", "", "Nameserver")

	domainCreateCmd.Flags().StringVarP(&domainCreateNsIp1, "nameserverIp1", "", "", "IP address for nameserver")
	domainCreateCmd.Flags().StringVarP(&domainCreateNsIp2, "nameserverIp2", "", "", "IP address for nameserver")
	domainCreateCmd.Flags().StringVarP(&domainCreateNsIp3, "nameserverIp3", "", "", "IP address for nameserver")
	domainCreateCmd.Flags().StringVarP(&domainCreateNsIp4, "nameserverIp4", "", "", "IP address for nameserver")

	domainCreateCmd.Flags().StringVarP(&domainCreateNsIpv61, "nameserverIpv61", "", "", "IPv6 address for nameserver")
	domainCreateCmd.Flags().StringVarP(&domainCreateNsIpv62, "nameserverIpv62", "", "", "IPv6 address for nameserver")
	domainCreateCmd.Flags().StringVarP(&domainCreateNsIpv63, "nameserverIpv63", "", "", "IPv6 address for nameserver")
	domainCreateCmd.Flags().StringVarP(&domainCreateNsIpv64, "nameserverIpv64", "", "", "IPv6 address for nameserver")

	domainCreateCmd.Flags().IntVarP(&domainCreateTtl, "ttl", "", 28800, "Time to live: amount of time (in seconds) the DNS-records stay in the cache")
	domainCreateCmd.Flags().StringVarP(&domainCreateEpCode, "eppCode", "", "", "eppCode")
	domainCreateCmd.Flags().BoolVarP(&domainCreateHandleDns, "handleDNS", "", false, "should dns be handled by lvl27")
	domainCreateCmd.Flags().StringVarP(&domainCreateExtraFields, "extra fields", "", "", "extra fields (json, non-editable)")

	domainCreateCmd.Flags().IntVarP(&domainCreateContactOnSite, "domaincontactOnsite", "", 0, "the unique id of a domaincontact with type onsite")
	domainCreateCmd.Flags().StringVarP(&domainCreateAutoRecordTemplate, "autorecordTemplate", "", "", "AutorecordTemplate")
	domainCreateCmd.Flags().BoolVarP(&domainCreateAutoRecordTemplateRep, "autorecordTemplateReplace", "", false, "autorecordTemplate replace")
	
	domainCreateCmd.Flags().IntVarP(&domainCreateDomainProvider, "domainProvider", "", 0, "The id of a domain provider (admin only)")
	domainCreateCmd.Flags().StringVarP(&domainCreateExternalCreated, "dtExternallCreated", "", "", "Creation timestamp (admin only)")
	domainCreateCmd.Flags().StringVarP(&domainCreateExternalExpires, "dtExternallExpires", "", "", "Expire date timestamp (admin only)")
	
	domainCreateCmd.Flags().StringVarP(&domainCreateConvertDomainRecords, "convertDomainrecords", "", "", "Domainrecord json (admin only)")
	domainCreateCmd.Flags().StringVarP(&domainCreateAutoTeams, "autoTeams", "", "", "a csv list of team id's")
	domainCreateCmd.Flags().StringVarP(&domainCreateExternalInfo, "externalInfo", "", "", "Required when billableItemInfo for an organisation exist in db")
	


	

	// ----------------- RECORDS ------------------------
	domainCmd.AddCommand(domainRecordCmd)

	// Record list
	domainRecordCmd.AddCommand(domainRecordListCmd)

	// Record create
	flags := domainRecordCreateCmd.Flags() 
	flags.StringVarP(&domainRecordCreateType, "type", "t", "", "Type of the domain record")
	flags.StringVarP(&domainRecordCreateName, "name", "n", "", "Name of the domain record")
	flags.StringVarP(&domainRecordCreateContent, "content", "c", "", "Content of the domain record")
	flags.IntVarP(&domainRecordCreatePriority, "priority", "p", 0, "Priority of the domain record")
	domainRecordCreateCmd.MarkFlagRequired("type")
	domainRecordCreateCmd.MarkFlagRequired("content")
	domainRecordCmd.AddCommand(domainRecordCreateCmd)

	// Record delete
	domainRecordCmd.AddCommand(domainRecordDeleteCmd)
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
// required flag vars
var domainCreateType,domainCreateLicensee,domainCreateOrganisation int
var domainCreateName string
// non-required flag vars
var domainCreateNs1,domainCreateNs2,domainCreateNs3,domainCreateNs4 string
var domainCreateNsIp1,domainCreateNsIp2,domainCreateNsIp3,domainCreateNsIp4 string
var domainCreateNsIpv61,domainCreateNsIpv62,domainCreateNsIpv63,domainCreateNsIpv64 string
var domainCreateTtl,domainCreateContactOnSite,domainCreateDomainProvider int
var domainCreateEpCode, domainCreateAutoRecordTemplate string
var domainCreateHandleDns,domainCreateAutoRecordTemplateRep bool
var domainCreateExtraFields,domainCreateExternalCreated,domainCreateExternalExpires string
var domainCreateConvertDomainRecords, domainCreateAutoTeams,domainCreateExternalInfo string


var domainCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new domain",

	Run: func(ccmd *cobra.Command, args []string) {
		Level27Client.DomainCreate(args, types.DomainRequest{
			Name: domainCreateName,
			Action: "none",
			Domaintype: domainCreateType,
			Domaincontactlicensee: domainCreateLicensee,
			Organisation: domainCreateOrganisation,
		})
	},
}

// ----------------- RECORDS ------------------------

var domainRecordCmd = &cobra.Command{
	Use:   "record",
	Short: "Commands for managing domain records",
}


// GET DOMAIN/RECORDS
var domainRecordListCmd = &cobra.Command{
	Use: "list [domain]",
	Short: "Get a list of all records configured for a domain",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		records := Level27Client.DomainRecords(args[0])
		
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tTYPE\tNAME\tCONTENT\t")

		for _, record := range records {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t\n", record.ID, record.Type, record.Name, record.Content)
		}
	
		w.Flush()
	},
}

// CREATE DOMAIN/RECORD
var domainRecordCreateType string
var domainRecordCreateName string
var domainRecordCreateContent string
var domainRecordCreatePriority int

var domainRecordCreateCmd = &cobra.Command{
	Use: "create [domain]",
	Short: "Create a new record for a domain",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalln("Not a valid domain ID!")
		}

		Level27Client.DomainRecordCreate(id, types.DomainRecordRequest{
			Name: domainRecordCreateName,
			Type: domainRecordCreateType,
			Priority: domainRecordCreatePriority,
			Content: domainRecordCreateContent,
		})
	},
}


// DELETE DOMAIN/RECORD
var domainRecordDeleteCmd = &cobra.Command{
	Use: "delete [domain] [record]",
	Short: "Delete a record for a domain",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		domainId, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalln("Not a valid domain ID!")
		}

		recordId, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalln("Not a valid domain ID!")
		}

		Level27Client.DomainRecordDelete(domainId, recordId)
	},
}






