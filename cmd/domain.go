package cmd

import (
	"fmt"
	"log"
	"strconv"

	"bitbucket.org/level27/lvl/types"
	"github.com/spf13/cobra"
)

// MAIN COMMAND
var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Commands for managing domains",
}

func init() {

	// ---------------------------------------------------- DOMAINS -------------------------------------------------------
	RootCmd.AddCommand(domainCmd)

	// Get (list of all domains)
	domainCmd.AddCommand(domainGetCmd)

	addCommonGetFlags(domainGetCmd)

	// Get details from a specific domain
	domainCmd.AddCommand(domainDescribeCmd)

	// Delete (single domain)
	domainCmd.AddCommand(domainDeleteCmd)

	// Create (single domain)
	domainCmd.AddCommand(domainCreateCmd)
	domainCreateCmd.Flags().StringVarP(&domainCreateAction, "action", "a", "", "Specify the action you want to commit")
	domainCreateCmd.Flags().StringVarP(&domainCreateExternalInfo, "externalInfo", "", "", "Required when billableItemInfo for an organisation exist in db")
	addDomainCommonPostFlags(domainCreateCmd)
	//Required flags
	domainCreateCmd.MarkFlagRequired("name")
	domainCreateCmd.MarkFlagRequired("type")
	domainCreateCmd.MarkFlagRequired("licensee")
	domainCreateCmd.MarkFlagRequired("organisation")

	// TRANSFER (single domain)
	domainCmd.AddCommand(domainTransferCmd)
	addDomainCommonPostFlags(domainTransferCmd)
	// required flags
	domainTransferCmd.MarkFlagRequired("name")
	domainTransferCmd.MarkFlagRequired("type")
	domainTransferCmd.MarkFlagRequired("licensee")
	domainTransferCmd.MarkFlagRequired("organisation")
	domainTransferCmd.MarkFlagRequired("eppCode")

	// INTERNAL TRANSFER
	domainCmd.AddCommand(domainInternalTransferCmd)
	addDomainCommonPostFlags(domainInternalTransferCmd)

	// UPDATE (single domain)
	domainCmd.AddCommand(domainUpdateCmd)
	addDomainCommonPostFlags(domainUpdateCmd)
	//required flags
	domainUpdateCmd.MarkFlagRequired("name")
	domainUpdateCmd.MarkFlagRequired("type")
	domainUpdateCmd.MarkFlagRequired("licensee")
	domainUpdateCmd.MarkFlagRequired("organisation")

	// ------------------------------------------------- RECORDS ---------------------------------------------------------
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

	// Record update
	flags = domainRecordUpdateCmd.Flags()
	flags.StringVarP(&domainRecordUpdateName, "name", "n", "", "Name of the domain record")
	flags.StringVarP(&domainRecordUpdateContent, "content", "c", "", "Content of the domain record")
	flags.IntVarP(&domainRecordUpdatePriority, "priority", "p", 0, "Priority of the domain record")
	domainRecordCmd.AddCommand(domainRecordUpdateCmd)

	// Record delete
	domainRecordCmd.AddCommand(domainRecordDeleteCmd)

	// --------------------------------------------------- ACCESS --------------------------------------------------------
	domainCmd.AddCommand(domainAccessCmd)

	// ADD ACCESS
	domainAccessCmd.AddCommand(domainAccessAddCmd)

	flags = domainAccessAddCmd.Flags()
	flags.IntVarP(&domainAccessAddOrganisation, "organisation", "", 0, "The unique identifier of an organisation")
	domainAccessAddCmd.MarkFlagRequired("organisation")

	// REMOVE ACCESS
	domainAccessCmd.AddCommand(domainAccessRemoveCmd)
	flags = domainAccessRemoveCmd.Flags()
	flags.IntVarP(&domainAccessAddOrganisation, "organisation", "", 0, "The unique identifier of an organisation")
	domainAccessRemoveCmd.MarkFlagRequired("organisation")

	// --------------------------------------------------- NOTIFICATIONS --------------------------------------------------------
	domainCmd.AddCommand(domainNotificationCmd)

	// CREATE NOTIFICATION
	domainNotificationCmd.AddCommand(domainNotificationsCreateCmd)
	flags = domainNotificationsCreateCmd.Flags()
	flags.StringVarP(&domainNotificationPostType, "type", "t", "", "The notification type")
	flags.StringVarP(&domainNotificationPostGroup, "group", "g", "", "The notification group")
	flags.StringVarP(&domainNotificationPostParams, "params", "p", "", "Additional parameters (json)")
	flags.SortFlags = false
	domainNotificationCmd.MarkFlagRequired("type")
	domainNotificationCmd.MarkFlagRequired("group")

}

// --------------------------------------------------- DOMAINS --------------------------------------------------------
//GET LIST OF ALL DOMAINS [lvl domain get]
var domainGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a list of all current domains",
	Run: func(ccmd *cobra.Command, args []string) {
		outputFormatTable(
			getDomains(args),
			[]string{"ID", "NAME", "STATUS"},
			[]string{"ID", "Fullname", "Status"})
	},
}

func getDomains(ids []string) []types.Domain {
	c := Level27Client
	if len(ids) == 0 {
		return c.Domains(optFilter, optNumber)
	} else {
		domains := make([]types.Domain, len(ids))
		for idx, id := range ids {
			domains[idx] = c.Domain("GET", id, nil)
		}
		return domains
	}
}

// DESCRIBE DOMAIN (get detailed info from specific domain) - [lvl domain describe <id>]
var domainDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Get detailed info about a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domainID := args[0]
		domain := Level27Client.Domain("GET", domainID, nil)

		outputFormatTemplate(domain, "templates/domain.tmpl")
	},
}

// DELETE DOMAIN [lvl domain delete <id>]
var domainDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a domain",
	Long:  "use LVL DOMAIN DELETE <ID or ID's>. You can give multiple ID's to this command by seperating them trough whitespaces.",
	Run: func(cmd *cobra.Command, args []string) {
		Level27Client.DomainDelete(args)
	},
}

//flag vars needed for all post or put requests on Domain level [Domains/]
var domainCreateType, domainCreateLicensee, domainCreateOrganisation int
var domainCreateName string
var domainCreateNs1, domainCreateNs2, domainCreateNs3, domainCreateNs4 string
var domainCreateNsIp1, domainCreateNsIp2, domainCreateNsIp3, domainCreateNsIp4 string
var domainCreateNsIpv61, domainCreateNsIpv62, domainCreateNsIpv63, domainCreateNsIpv64 string
var domainCreateTtl, domainCreateDomainProvider int
var domainCreateEppCode, domainCreateAutoRecordTemplate string
var domainCreateHandleDns, domainCreateAutoRecordTemplateRep bool
var domainCreateExtraFields, domainCreateExternalCreated, domainCreateExternalExpires string
var domainCreateConvertDomainRecords, domainCreateAutoTeams, domainCreateExternalInfo, domainCreateAction string
var domainCreateContactOnSite int

// common functions for managing domains
// change given flag data into request data to put or post
func getDomainRequestData() types.DomainRequest {
	requestData := types.DomainRequest{
		Name:        domainCreateName,
		NameServer1: &domainCreateNs1,
		NameServer2: domainCreateNs2,
		NameServer3: domainCreateNs3,
		NameServer4: domainCreateNs4,

		NameServer1Ip: domainCreateNsIp1,
		NameServer2Ip: domainCreateNsIp2,
		NameServer3Ip: domainCreateNsIp3,
		NameServer4Ip: domainCreateNsIp4,

		NameServer1Ipv6: domainCreateNsIpv61,
		NameServer2Ipv6: domainCreateNsIpv62,
		NameServer3Ipv6: domainCreateNsIpv63,
		NameServer4Ipv6: domainCreateNsIpv64,

		TTL:                       domainCreateTtl,
		Action:                    domainCreateAction,
		EppCode:                   domainCreateEppCode,
		Handledns:                 domainCreateHandleDns,
		ExtraFields:               domainCreateExtraFields,
		Domaintype:                domainCreateType,
		Domaincontactlicensee:     domainCreateLicensee,
		DomainContactOnSite:       &domainCreateContactOnSite,
		Organisation:              domainCreateOrganisation,
		AutoRecordTemplate:        domainCreateAutoRecordTemplate,
		AutoRecordTemplateReplace: domainCreateAutoRecordTemplateRep,
		//DomainProvider:            &domainCreateDomainProvider,
		// DtExternalCreated:         domainCreateExternalCreated,
		// DtExternalExpires:         domainCreateExternalExpires,
		// ConvertDomainRecords:      domainCreateConvertDomainRecords,
		AutoTeams:    domainCreateAutoTeams,
		ExternalInfo: domainCreateExternalInfo,
	}

	if *requestData.DomainContactOnSite == 0 {
		requestData.DomainContactOnSite = nil
	}

	return requestData
}

// get all possible domain extensions en their ID
func getDomainExtensions() {

	res := Level27Client.Extension()

	fmt.Println(res[0])
}

// CREATE DOMAIN [lvl domain create (action:create/none)]
var domainCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new domain",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		getDomainExtensions()

		// requestData := getDomainRequestData()

		// if cmd.Flags().Changed("action") {

		// 	if requestData.Action == "create" {
		// 		Level27Client.DomainCreate(args, requestData)

		// 	} else if requestData.Action == "none" {
		// 		Level27Client.DomainCreate(args, requestData)
		// 	} else {
		// 		log.Printf("given action: '%v' is not recognized.", requestData.Action)
		// 	}
		// } else {
		// 	Level27Client.DomainCreate(args, requestData)
		// }

	},
}

// TRANSFER DOMAIN
var domainTransferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Command for transfering a domain",
	Run: func(cmd *cobra.Command, args []string) {

		requestData := getDomainRequestData()
		Level27Client.DomainTransfer(args, requestData)
	},
}

//INTERNAL TRANSFER
var domainInternalTransferCmd = &cobra.Command{
	Use:   "internaltransfer",
	Short: "Internal transfer (available only for dnsbe domains)",
	Run: func(cmd *cobra.Command, args []string) {

		requestData := getDomainRequestData()
		Level27Client.DomainTransfer(args, requestData)
	},
}

// UPDATE DOMAIN
var domainUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Command for updating an existing domain",

	Run: func(cmd *cobra.Command, args []string) {
		requestData := types.DomainUpdateRequest{
			Name:        domainCreateName,
			NameServer1: &domainCreateNs1,
			NameServer2: domainCreateNs2,
			NameServer3: domainCreateNs3,
			NameServer4: domainCreateNs4,

			NameServer1Ip: domainCreateNsIp1,
			NameServer2Ip: domainCreateNsIp2,
			NameServer3Ip: domainCreateNsIp3,
			NameServer4Ip: domainCreateNsIp4,

			NameServer1Ipv6: domainCreateNsIpv61,
			NameServer2Ipv6: domainCreateNsIpv62,
			NameServer3Ipv6: domainCreateNsIpv63,
			NameServer4Ipv6: domainCreateNsIpv64,

			TTL:                       domainCreateTtl,
			Action:                    domainCreateAction,
			EppCode:                   domainCreateEppCode,
			Handledns:                 domainCreateHandleDns,
			ExtraFields:               domainCreateExtraFields,
			Domaintype:                domainCreateType,
			Domaincontactlicensee:     domainCreateLicensee,
			DomainContactOnSite:       &domainCreateContactOnSite,
			Organisation:              domainCreateOrganisation,
			AutoRecordTemplate:        domainCreateAutoRecordTemplate,
			AutoRecordTemplateReplace: domainCreateAutoRecordTemplateRep,
			//DomainProvider:            &domainCreateDomainProvider,
			// DtExternalCreated:         domainCreateExternalCreated,
			// DtExternalExpires:         domainCreateExternalExpires,
			// ConvertDomainRecords:      domainCreateConvertDomainRecords,
			AutoTeams: domainCreateAutoTeams,
		}

		if *requestData.DomainContactOnSite == 0 {
			requestData.DomainContactOnSite = nil
		}

		Level27Client.DomainUpdate(args, requestData)
	},
}

// --------------------------------------------------- RECORDS --------------------------------------------------------

var domainRecordCmd = &cobra.Command{
	Use:   "record",
	Short: "Commands for managing domain records",
}

// GET DOMAIN/RECORDS
var domainRecordListCmd = &cobra.Command{
	Use:   "list [domain]",
	Short: "Get a list of all records configured for a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		records := Level27Client.DomainRecords(args[0])

		outputFormatTable(records, []string{"ID", "TYPE", "NAME", "CONTENT"}, []string{"ID", "Type", "Name", "Content"})
	},
}

// CREATE DOMAIN/RECORD
var domainRecordCreateType string
var domainRecordCreateName string
var domainRecordCreateContent string
var domainRecordCreatePriority int

var domainRecordCreateCmd = &cobra.Command{
	Use:   "create [domain]",
	Short: "Create a new record for a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalln("Not a valid domain ID!")
		}

		Level27Client.DomainRecordCreate(id, types.DomainRecordRequest{
			Name:     domainRecordCreateName,
			Type:     domainRecordCreateType,
			Priority: domainRecordCreatePriority,
			Content:  domainRecordCreateContent,
		})
	},
}

// DELETE DOMAIN/RECORD
var domainRecordDeleteCmd = &cobra.Command{
	Use:   "delete [domain] [record]",
	Short: "Delete a record for a domain",
	Args:  cobra.ExactArgs(2),
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

var domainRecordUpdateName string
var domainRecordUpdateContent string
var domainRecordUpdatePriority int

var domainRecordUpdateCmd = &cobra.Command{
	Use:   "update [domain] [record]",
	Short: "Update a record for a domain",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		domainId, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalln("Not a valid domain ID!")
		}

		recordId, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalln("Not a valid domain ID!")
		}

		// Merge data with existing so we don't bulldoze anything.
		data := Level27Client.DomainRecord(domainId, recordId)
		request := types.DomainRecordRequest{
			Type:     data.Type,
			Name:     data.Name,
			Content:  data.Content,
			Priority: data.Priority,
		}

		if cmd.Flags().Changed("name") {
			request.Name = domainRecordUpdateName
		}

		if cmd.Flags().Changed("content") {
			request.Content = domainRecordUpdateContent
		}

		if cmd.Flags().Changed("priority") {
			request.Priority = domainRecordUpdatePriority
		}

		Level27Client.DomainRecordUpdate(domainId, recordId, request)
	},
}

// --------------------------------------------------- ACCESS --------------------------------------------------------
var domainAccessCmd = &cobra.Command{
	Use:   "access",
	Short: "Commands for managing the access of a domain",
}

// ADD ACCESS TO A DOMAIN
var domainAccessAddOrganisation int

var domainAccessAddCmd = &cobra.Command{
	Use:   "add [domain] [flags]",
	Short: "Add organisation access to a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalln("Not a valid domain ID!")
		}

		Level27Client.DomainAccesAdd(id, types.DomainAccessRequest{
			Organisation: domainAccessAddOrganisation,
		})
	},
}

// REMOVE ACCESS FROM DOMAIN
var domainAccessRemoveCmd = &cobra.Command{
	Use:   "delete [domain] [flags]",
	Short: "Remove organisation acces from a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalln("Not a valid domain ID!")
		}

		var orgId int

		if cmd.Flags().Changed("organisation") {
			value := cmd.Flag("organisation").Value.String()
			orgId, err = strconv.Atoi(value)
			if err != nil {
				log.Fatal("no valid organisation ID")
			}
			Level27Client.DomainAccesRemove(id, orgId)
		}

	},
}


// --------------------------------------------------- NOTIFICATIONS --------------------------------------------------------
// MAIN COMMAND
var domainNotificationCmd = &cobra.Command{
	Use:   "notification",
	Short: "Commands for managing domain notifications",
}

// CREATE NOTIFICATION
var domainNotificationPostType, domainNotificationPostGroup, domainNotificationPostParams string

var domainNotificationsCreateCmd = &cobra.Command{
	Use: "create [domain] [flags]",
	Short: "Send a notification for a domain",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id , err := strconv.Atoi(args[0])
		if err != nil{
			log.Fatal("no valid domain ID")
		}

		Level27Client.DomainNotificationAdd(id, types.DomainNotificationPostRequest{
			Type: domainNotificationPostType,
			Group: domainNotificationPostGroup,
			Params: domainNotificationPostParams,
		})

	},
}
