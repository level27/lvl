package cmd

import (
	"fmt"

	"log"
	"strconv"
	"strings"

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
	domainCreateCmd.MarkFlagRequired("licensee")
	domainCreateCmd.MarkFlagRequired("organisation")

	// TRANSFER (single domain)
	domainCmd.AddCommand(domainTransferCmd)
	addDomainCommonPostFlags(domainTransferCmd)
	// required flags
	domainTransferCmd.MarkFlagRequired("name")
	domainTransferCmd.MarkFlagRequired("licensee")
	domainTransferCmd.MarkFlagRequired("organisation")
	domainTransferCmd.MarkFlagRequired("eppCode")

	// INTERNAL TRANSFER
	domainCmd.AddCommand(domainInternalTransferCmd)
	addDomainCommonPostFlags(domainInternalTransferCmd)

	// UPDATE (single domain)
	domainCmd.AddCommand(domainUpdateCmd)
	settingString(domainUpdateCmd, domainUpdateSettings, "nameserver1", "")
	settingString(domainUpdateCmd, domainUpdateSettings, "nameserver2", "")
	settingString(domainUpdateCmd, domainUpdateSettings, "nameserver3", "")
	settingString(domainUpdateCmd, domainUpdateSettings, "nameserverIp1", "")
	settingString(domainUpdateCmd, domainUpdateSettings, "nameserverIp2", "")
	settingString(domainUpdateCmd, domainUpdateSettings, "nameserverIp3", "")
	settingString(domainUpdateCmd, domainUpdateSettings, "nameserverIpv61", "")
	settingString(domainUpdateCmd, domainUpdateSettings, "nameserverIpv62", "")
	settingString(domainUpdateCmd, domainUpdateSettings, "nameserverIpv63", "")
	settingInt(domainUpdateCmd, domainUpdateSettings, "ttl", "")
	settingBool(domainUpdateCmd, domainUpdateSettings, "handleDns", "")
	settingInt(domainUpdateCmd, domainUpdateSettings, "domaincontactLicensee", "")
	settingInt(domainUpdateCmd, domainUpdateSettings, "domaincontactOnSite", "")
	settingInt(domainUpdateCmd, domainUpdateSettings, "organisation", "")

	// ------------------------------------------------- RECORDS ---------------------------------------------------------
	domainCmd.AddCommand(domainRecordCmd)

	// Record list
	domainRecordCmd.AddCommand(domainRecordGetCmd)
	addCommonGetFlags(domainRecordGetCmd)
	domainRecordGetCmd.Flags().StringVarP(&recordGetType, "type", "t", "", "Type of records to filter")

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

	/*
		// --------------------------------------------------- NOTIFICATIONS --------------------------------------------------------
		domainCmd.AddCommand(domainNotificationCmd)

		// CREATE NOTIFICATION
		domainNotificationCmd.AddCommand(domainNotificationsCreateCmd)
		flags = domainNotificationsCreateCmd.Flags()
		flags.StringVarP(&domainNotificationPostType, "type", "t", "", "The notification type")
		flags.StringVarP(&domainNotificationPostGroup, "group", "g", "", "The notification group")
		flags.StringVarP(&domainNotificationPostParams, "params", "p", "", "Additional parameters (json)")
		flags.SortFlags = false
		domainNotificationsCreateCmd.MarkFlagRequired("type")
		domainNotificationsCreateCmd.MarkFlagRequired("group")

		// GET NOTIFICATIONS
		var notificationsOrderBy string
		domainNotificationCmd.AddCommand(domainNotificationsGetCmd)
		flags = domainNotificationsGetCmd.Flags()
		flags.StringVarP(&notificationsOrderBy, "orderby", "", "", "The field you want to order the results on")
		flags.SortFlags = false
		addCommonGetFlags(domainNotificationsGetCmd)
	*/
	// --------------------------------------------------- BILLABLEITEMS --------------------------------------------------------
	domainCmd.AddCommand(domainBillableItemCmd)

	// CREATE BILLABLEITEM (turn invoicing on)
	domainBillableItemCmd.AddCommand(domainBillCreateCmd)
	flags = domainBillCreateCmd.Flags()
	flags.StringVarP(&externalInfo, "externalinfo", "e", "", "ExternalInfo (required when billableitemInfo entities for an Organisation exist in db)")

	// DELETE BILLABLEITEM (turn invoicing off)
	domainBillableItemCmd.AddCommand(domainBillDeleteCmd)

	// --------------------------------------------------- AVAILABILITY/CHECK --------------------------------------------------------
	// CHECK
	domainCmd.AddCommand(domainCheckCmd)

	// --------------------------------------------------- JOB HISTORY --------------------------------------------------------
	domainCmd.AddCommand(domainJobHistoryCmd)

	domainCmd.AddCommand(domainRootJobHistoryCmd)

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
	Use:   "delete [domainId]",
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
var domainCreateTtl int
var domainCreateEppCode, domainCreateAutoRecordTemplate string
var domainCreateHandleDns, domainCreateAutoRecordTemplateRep bool
var domainCreateExtraFields string
var domainCreateAutoTeams, domainCreateExternalInfo, domainCreateAction string
var domainCreateContactOnSite int

// var domainCreateConvertDomainRecords, domainCreateExternalCreated, domainCreateExternalExpires string
// var  domainCreateDomainProvider int

// common functions for managing domains
// change given flag data into request data to put or post
func getDomainRequestData() types.DomainRequest {
	requestData := types.DomainRequest{
		Name:          domainCreateName,
		NameServer1:   &domainCreateNs1,
		NameServer2:   domainCreateNs2,
		NameServer3:   domainCreateNs3,
		NameServer4:   domainCreateNs4,
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

	if requestData.Domaintype == 0 {
		name, extension, domainType := getDomainTypeForDomain(requestData.Name)
		if domainType == 0 {
			log.Fatalf("Invalid domain extension: '%s'", extension)
		}

		requestData.Domaintype = domainType
		requestData.Name = name
	}

	return requestData
}

// Splits a domain name into its name and extension respectively.
func splitDomainName(domain string) (string, string) {
	idx := strings.IndexByte(domain, '.')
	extension := domain[idx+1:]
	name := domain[:idx]

	return name, extension
}

// Gets the domain type extension for a full domain name.
func getDomainTypeForDomain(domain string) (string, string, int) {
	name, extension := splitDomainName(domain)
	res := Level27Client.Extension()

	for _, provider := range res {
		for _, domainType := range provider.Domaintypes {
			if domainType.Extension == extension {
				return name, extension, domainType.ID
			}
		}
	}

	return name, extension, 0
}

// CREATE DOMAIN [lvl domain create (action:create/none)]
var domainCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new domain",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		requestData := getDomainRequestData()

		if cmd.Flags().Changed("action") {

			if requestData.Action == "create" {
				Level27Client.DomainCreate(args, requestData)

			} else if requestData.Action == "none" {
				Level27Client.DomainCreate(args, requestData)
			} else {
				log.Printf("given action: '%v' is not recognized.", requestData.Action)
			}
		} else {
			Level27Client.DomainCreate(args, requestData)
		}

	},
}

// TRANSFER DOMAIN
var domainTransferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer a domain",
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
var domainUpdateSettings map[string]interface{} = make(map[string]interface{})

// UPDATE DOMAIN
var domainUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Command for updating an existing domain",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		domainId, err := convertStringToId(args[0])
		cobra.CheckErr(err)

		if len(domainUpdateSettings) == 0 {
			fmt.Println("No options specified!")
		}

		Level27Client.DomainUpdate(domainId, domainUpdateSettings)
	},
}

// --------------------------------------------------- RECORDS --------------------------------------------------------

var domainRecordCmd = &cobra.Command{
	Use:   "record",
	Short: "Manage domain records",
}

var recordGetType string

// GET DOMAIN/RECORDS
var domainRecordGetCmd = &cobra.Command{
	Use:   "get [domain]",
	Short: "Get a list of all records configured for a domain",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domainId, err := convertStringToId(args[0])
		cobra.CheckErr(err)
		recordIds, err := convertStringsToIds(args[1:])
		cobra.CheckErr(err)

		records := getDomainRecords(domainId, recordIds)

		outputFormatTable(records, []string{"ID", "TYPE", "NAME", "CONTENT"}, []string{"ID", "Type", "Name", "Content"})
	},
}

func getDomainRecords(domainId int, ids []int) []types.DomainRecord {
	c := Level27Client
	if len(ids) == 0 {
		return c.DomainRecords(domainId, recordGetType, optNumber, optFilter)
	} else {
		domains := make([]types.DomainRecord, len(ids))
		for idx, id := range ids {
			domains[idx] = c.DomainRecord(domainId, id)
		}
		return domains
	}
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
	Short: "Manage the access of a domain",
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
	Short: "Remove organisation access from a domain",
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
/*
// MAIN COMMAND
var domainNotificationCmd = &cobra.Command{
	Use:   "notifications",
	Short: "Manage domain notifications",
}

// CREATE NOTIFICATION
var domainNotificationPostType, domainNotificationPostGroup, domainNotificationPostParams string

var domainNotificationsCreateCmd = &cobra.Command{
	Use:   "create [domain] [flags]",
	Short: "Send a notification for a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal("no valid domain ID")
		}

		Level27Client.DomainNotificationAdd(id, types.DomainNotificationPostRequest{
			Type:   domainNotificationPostType,
			Group:  domainNotificationPostGroup,
			Params: domainNotificationPostParams,
		})

	},
}

// GET NOTIFICATIONS
var domainNotificationsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get a list of all notifications from a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal("no valid domain ID")
		}

		notifications := Level27Client.DomainNotificationGet(id)
		fmt.Print(notifications)

	},
}
*/
// --------------------------------------------------- BILLABLEITEMS --------------------------------------------------------
// MAIN COMMAND
var domainBillableItemCmd = &cobra.Command{
	Use:   "billing",
	Short: "Manage domain's invoicing (BillableItem)",
}

// CREATE A BILLABLEITEM / TURN ON BILLING(ADMIN ONLY)
var externalInfo string

var domainBillCreateCmd = &cobra.Command{
	Use:   "on [domain] [flags]",
	Short: "Turn on billing for a domain (admin only)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal("no valid domain ID")
		}
		req := types.DomainBillPostRequest{
			ExternalInfo: externalInfo,
		}

		Level27Client.DomainBillableItemCreate(id, req)

	},
}

//DELETE BILLABLEITEM/ TURN OF BILLING

var domainBillDeleteCmd = &cobra.Command{
	Use:   "off [domainID]",
	Short: "Turn off the billing for domain (admin only)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal("no valid domain ID")
		}

		Level27Client.DomainBillableItemDelete(id)

	},
}

// ---------------------------------------------- CHECK / AVAILABILITY ------------------------------------------------
var domainCheckCmd = &cobra.Command{
	Use:   "check [domain name]",
	Short: "Check availability of a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		name, extension := splitDomainName(domain)

		status := Level27Client.DomainCheck(name, extension)

		outputFormatTemplate(status, "templates/domainCheck.tmpl")
	},
}

// ---------------------------------------------- JOB HISTORY DOMAINS ------------------------------------------------

var domainJobHistoryCmd = &cobra.Command{
	Use:   "jobs [domainId]",
	Short: "Manage the job history for a domain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal("no valid domain ID")
		}
		history := Level27Client.DomainJobHistoryGet(id)

		outputFormatTable(history, []string{"ID", "STATUS", "MESSAGE"}, []string{"Id", "Status", "Message"})

	},
}

var domainRootJobHistoryCmd = &cobra.Command{
	Use:   "root",
	Short: "get detailed jobs",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal("no valid domain ID")
		}
		fmt.Print(Level27Client.DomainJobHistoryRootGet(id))
	},
}
