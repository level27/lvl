package cmd

import (
	"fmt"

	"strconv"
	"strings"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

// MAIN COMMAND
var domainCmd = &cobra.Command{
	Use:     "domain",
	Short:   "Commands for managing domains",
	Example: "lvl domain get -f example.be",
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
	addWaitFlag(domainDeleteCmd)
	addDeleteConfirmFlag(domainDeleteCmd)

	// Create (single domain)
	domainCmd.AddCommand(domainCreateCmd)
	addWaitFlag(domainCreateCmd)
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
	addAccessCmds(domainCmd, "domains", resolveDomain)

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
	addBillingCmds(domainCmd, "domains", resolveDomain)

	// --------------------------------------------------- AVAILABILITY/CHECK --------------------------------------------------------
	// CHECK
	domainCmd.AddCommand(domainCheckCmd)

	// --------------------------------------------------- JOB HISTORY --------------------------------------------------------
	addJobCmds(domainCmd, "domain", resolveDomain)

	// INTEGRITY CHECKS
	addIntegrityCheckCmds(domainCmd, "domains", resolveDomain)
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

// common date used for Post operations at /Domains
func addDomainCommonPostFlags(cmd *cobra.Command) {
	command := cmd.Flags()

	command.StringVarP(&domainCreateName, "name", "n", "", "the name of the domain (REQUIRED)")
	command.IntVarP(&domainCreateType, "type", "t", 0, "the type of the domain")
	command.MarkHidden("type")
	command.IntVarP(&domainCreateLicensee, "licensee", "l", 0, "The unique identifier of a domaincontact with type licensee (REQUIRED)")
	command.IntVarP(&domainCreateOrganisation, "organisation", "", 0, "the organisation of the domain (REQUIRED)")

	command.StringVarP(&domainCreateNs1, "nameserver1", "", "", "Nameserver")
	command.StringVarP(&domainCreateNs2, "nameserver2", "", "", "Nameserver")
	command.StringVarP(&domainCreateNs3, "nameserver3", "", "", "Nameserver")
	command.StringVarP(&domainCreateNs4, "nameserver4", "", "", "Nameserver")

	command.StringVarP(&domainCreateNsIp1, "nameserverIp1", "", "", "IP address for nameserver")
	command.StringVarP(&domainCreateNsIp2, "nameserverIp2", "", "", "IP address for nameserver")
	command.StringVarP(&domainCreateNsIp3, "nameserverIp3", "", "", "IP address for nameserver")
	command.StringVarP(&domainCreateNsIp4, "nameserverIp4", "", "", "IP address for nameserver")

	command.StringVarP(&domainCreateNsIpv61, "nameserverIpv61", "", "", "IPv6 address for nameserver")
	command.StringVarP(&domainCreateNsIpv62, "nameserverIpv62", "", "", "IPv6 address for nameserver")
	command.StringVarP(&domainCreateNsIpv63, "nameserverIpv63", "", "", "IPv6 address for nameserver")
	command.StringVarP(&domainCreateNsIpv64, "nameserverIpv64", "", "", "IPv6 address for nameserver")

	command.IntVarP(&domainCreateTtl, "ttl", "", 28800, "Time to live: amount of time (in seconds) the DNS-records stay in the cache")
	command.StringVarP(&domainCreateEppCode, "eppCode", "", "", "eppCode")
	command.BoolVarP(&domainCreateHandleDns, "handleDns", "", true, "should dns be handled by lvl27")
	command.StringVarP(&domainCreateExtraFields, "extra fields", "", "", "extra fields (json, non-editable)")

	command.IntVarP(&domainCreateContactOnSite, "domaincontactOnsite", "", 0, "the unique id of a domaincontact with type onsite")

	// command.StringVarP(&domainCreateAutoRecordTemplate, "autorecordTemplate", "", "", "AutorecordTemplate")
	// command.BoolVarP(&domainCreateAutoRecordTemplateRep, "autorecordTemplateReplace", "", false, "autorecordTemplate replace")
	//command.IntVarP(&domainCreateDomainProvider, "domainProvider", "", 0, "The id of a domain provider (admin only)")
	// command.StringVarP(&domainCreateExternalCreated, "dtExternallCreated", "", "", "Creation timestamp (admin only)")
	// command.StringVarP(&domainCreateExternalExpires, "dtExternallExpires", "", "", "Expire date timestamp (admin only)")
	// command.StringVarP(&domainCreateConvertDomainRecords, "convertDomainrecords", "", "", "Domainrecord json (admin only)")
	command.StringVarP(&domainCreateAutoTeams, "autoTeams", "", "", "a csv list of team id's")

	command.SortFlags = false
}

// Resolve an integer or name domain.
// If the domain is a name, a request is made to resolve the integer ID.
func resolveDomain(arg string) (int, error) {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.LookupDomain(arg)
	if err == nil {
		return 0, nil
	}

	res, err := resolveShared(
		options,
		arg,
		"domain",
		func(domain l27.Domain) string { return fmt.Sprintf("%s (%d)", domain.Name, domain.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, nil
}

// --------------------------------------------------- DOMAINS --------------------------------------------------------
//GET LIST OF ALL DOMAINS [lvl domain get]
var domainGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a list of all current domains",
	RunE: func(ccmd *cobra.Command, args []string) error {
		ids, err := convertStringsToIds(args)
		if err != nil {
			return err
		}

		options, err := getDomains(ids)
		if err != nil {
			return err
		}

		outputFormatTable(
			options,
			[]string{"ID", "NAME", "STATUS"},
			[]string{"ID", "Fullname", "Status"})

		return nil
	},
}

func getDomains(ids []int) ([]l27.Domain, error) {
	c := Level27Client
	if len(ids) == 0 {
		return c.Domains(optGetParameters)
	} else {
		domains := make([]l27.Domain, len(ids))
		for idx, id := range ids {
			var err error
			domains[idx], err = c.Domain(id)
			if err != nil {
				return nil, err
			}
		}

		return domains, nil
	}
}

// DESCRIBE DOMAIN (get detailed info from specific domain) - [lvl domain describe <id>]
var domainDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Get detailed info about a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domainID, err := resolveDomain(args[0])
		if err != nil {
			return err
		}

		domain, err := Level27Client.Domain(domainID)
		if err != nil {
			return err
		}

		domain.Jobs, err = Level27Client.EntityJobHistoryGet("domain", domainID)
		if err != nil {
			return err
		}

		for idx, j := range domain.Jobs {
			domain.Jobs[idx], err = Level27Client.JobHistoryRootGet(j.Id)
			if err != nil {
				return err
			}
		}

		outputFormatTemplate(domain, "templates/domain.tmpl")
		return nil
	},
}

// DELETE DOMAIN [lvl domain delete <id>]
var domainDeleteCmd = &cobra.Command{
	Use:   "delete [domainId]",
	Short: "Delete a domain",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		domainID, err := resolveDomain(args[0])
		if err != nil {
			return err
		}

		if !optDeleteConfirmed {
			domain, err := Level27Client.Domain(domainID)
			if err != nil {
				return err
			}

			if !confirmPrompt(fmt.Sprintf("Delete domain %s (%d)?", domain.Name, domain.ID)) {
				return nil
			}
		}

		err = Level27Client.DomainDelete(domainID)
		if err != nil {
			return err
		}

		if optWait {
			err = waitForDelete(
				func() (l27.Domain, error) { return Level27Client.Domain(domainID) },
				func(a l27.Domain) string { return a.Status },
				[]string{"deleting", "to_delete"},
			)

			if err != nil {
				return fmt.Errorf("waiting on domain status failed: %s", err.Error())
			}
		}

		return nil
	},
}

// common functions for managing domains
// change given flag data into request data to put or post
func getDomainRequestData() (l27.DomainRequest, error) {
	requestData := l27.DomainRequest{
		Name:          domainCreateName,
		NameServer1:   &domainCreateNs1,
		NameServer2:   &domainCreateNs2,
		NameServer3:   &domainCreateNs3,
		NameServer4:   &domainCreateNs4,
		NameServer1Ip: &domainCreateNsIp1,
		NameServer2Ip: &domainCreateNsIp2,
		NameServer3Ip: &domainCreateNsIp3,
		NameServer4Ip: &domainCreateNsIp4,

		NameServer1Ipv6: &domainCreateNsIpv61,
		NameServer2Ipv6: &domainCreateNsIpv62,
		NameServer3Ipv6: &domainCreateNsIpv63,
		NameServer4Ipv6: &domainCreateNsIpv64,

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
		ExternalInfo: &domainCreateExternalInfo,
	}

	if *requestData.DomainContactOnSite == 0 {
		requestData.DomainContactOnSite = nil
	}

	if requestData.Domaintype == 0 {
		name, extension, domainType, err := getDomainTypeForDomain(requestData.Name)
		if err != nil {
			return l27.DomainRequest{}, err
		}

		if domainType == 0 {
			return l27.DomainRequest{}, fmt.Errorf("invalid domain extension: '%s'", extension)
		}

		requestData.Domaintype = domainType
		requestData.Name = name
	}

	return requestData, nil
}

// Splits a domain name into its name and extension respectively.
func splitDomainName(domain string) (string, string) {
	idx := strings.IndexByte(domain, '.')
	extension := domain[idx+1:]
	name := domain[:idx]

	return name, extension
}

// Gets the domain type extension for a full domain name.
func getDomainTypeForDomain(domain string) (string, string, int, error) {
	name, extension := splitDomainName(domain)
	res, err := Level27Client.Extension()
	if err != nil {
		return "", "", 0, err
	}

	for _, provider := range res {
		for _, domainType := range provider.Domaintypes {
			if domainType.Extension == extension {
				return name, extension, domainType.ID, nil
			}
		}
	}

	return name, extension, 0, nil
}

// CREATE DOMAIN [lvl domain create (action:create/none)]
var domainCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new domain",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		requestData, err := getDomainRequestData()
		if err != nil {
			return err
		}

		domain, err := Level27Client.DomainCreate(requestData)
		if err != nil {
			return err
		}

		if optWait {
			domain, err = waitForStatus(
				func() (l27.Domain, error) { return Level27Client.Domain(domain.ID) },
				func(s l27.Domain) string { return s.Status },
				"ok",
				[]string{"to_create", "creating"},
			)

			if err != nil {
				return fmt.Errorf("waiting on domain status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(domain, "templates/entities/domain/create.tmpl")
		return nil
	},
}

// TRANSFER DOMAIN
var domainTransferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer a domain",
	RunE: func(cmd *cobra.Command, args []string) error {
		requestData, err := getDomainRequestData()
		if err != nil {
			return err
		}

		Level27Client.DomainTransfer(requestData)
		return nil
	},
}

//INTERNAL TRANSFER
var domainInternalTransferCmd = &cobra.Command{
	Use:   "internaltransfer",
	Short: "Internal transfer (available only for dnsbe domains)",
	RunE: func(cmd *cobra.Command, args []string) error {
		requestData, err := getDomainRequestData()
		if err != nil {
			return err
		}

		Level27Client.DomainTransfer(requestData)
		return nil
	},
}

// UPDATE DOMAIN
var domainUpdateSettings map[string]interface{} = make(map[string]interface{})
var domainUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Command for updating an existing domain",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		domainId, err := resolveDomain(args[0])
		if err != nil {
			return err
		}

		if len(domainUpdateSettings) == 0 {
			fmt.Println("No options specified!")
		}

		Level27Client.DomainUpdate(domainId, domainUpdateSettings)
		return nil
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
	RunE: func(cmd *cobra.Command, args []string) error {
		domainId, err := resolveDomain(args[0])
		if err != nil {
			return err
		}

		recordIds, err := convertStringsToIds(args[1:])
		if err != nil {
			return err
		}

		records, err := getDomainRecords(domainId, recordIds)
		if err != nil {
			return err
		}

		outputFormatTable(records, []string{"ID", "TYPE", "NAME", "CONTENT"}, []string{"ID", "Type", "Name", "Content"})
		return nil
	},
}

func getDomainRecords(domainId int, ids []int) ([]l27.DomainRecord, error) {
	c := Level27Client
	if len(ids) == 0 {
		return c.DomainRecords(domainId, recordGetType, optGetParameters)
	} else {
		domains := make([]l27.DomainRecord, len(ids))
		for idx, id := range ids {
			var err error
			domains[idx], err = c.DomainRecord(domainId, id)
			if err != nil {
				return nil, err
			}
		}

		return domains, nil
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
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := resolveDomain(args[0])
		if err != nil {
			return err
		}

		_, err = Level27Client.DomainRecordCreate(id, l27.DomainRecordRequest{
			Name:     domainRecordCreateName,
			Type:     domainRecordCreateType,
			Priority: domainRecordCreatePriority,
			Content:  domainRecordCreateContent,
		})

		return err
	},
}

// DELETE DOMAIN/RECORD
var domainRecordDeleteCmd = &cobra.Command{
	Use:   "delete [domain] [record]",
	Short: "Delete a record for a domain",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		//check for valid domain id
		domainId, err := resolveDomain(args[0])
		if err != nil {
			return err
		}

		recordId, err := checkSingleIntID(args[1], "record")
		if err != nil {
			return err
		}

		err = Level27Client.DomainRecordDelete(domainId, recordId)
		return err
	},
}

var domainRecordUpdateName string
var domainRecordUpdateContent string
var domainRecordUpdatePriority int

var domainRecordUpdateCmd = &cobra.Command{
	Use:   "update [domain] [record]",
	Short: "Update a record for a domain",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		domainId, err := resolveDomain(args[0])
		if err != nil {
			return err
		}

		recordId, err := checkSingleIntID(args[1], "record")
		if err != nil {
			return err
		}

		// Merge data with existing so we don't bulldoze anything.
		data, err := Level27Client.DomainRecord(domainId, recordId)
		if err != nil {
			return err
		}
		request := l27.DomainRecordRequest{
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
		return nil
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

		Level27Client.DomainNotificationAdd(id, DomainNotificationPostRequest{
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

// ---------------------------------------------- CHECK / AVAILABILITY ------------------------------------------------
var domainCheckCmd = &cobra.Command{
	Use:   "check [domain name]",
	Short: "Check availability of a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		name, extension := splitDomainName(domain)

		status, err := Level27Client.DomainCheck(name, extension)
		if err != nil {
			return err
		}

		outputFormatTemplate(status, "templates/domainCheck.tmpl")
		return nil
	},
}
