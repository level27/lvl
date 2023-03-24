package cmd

import (
	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

func init() {
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
	flags.Int32VarP(&domainRecordCreatePriority, "priority", "p", 0, "Priority of the domain record")
	domainRecordCreateCmd.MarkFlagRequired("type")
	domainRecordCreateCmd.MarkFlagRequired("content")
	domainRecordCmd.AddCommand(domainRecordCreateCmd)

	// Record update
	flags = domainRecordUpdateCmd.Flags()
	flags.StringVarP(&domainRecordUpdateName, "name", "n", "", "Name of the domain record")
	flags.StringVarP(&domainRecordUpdateContent, "content", "c", "", "Content of the domain record")
	flags.Int32VarP(&domainRecordUpdatePriority, "priority", "p", 0, "Priority of the domain record")
	domainRecordCmd.AddCommand(domainRecordUpdateCmd)

	// Record delete
	domainRecordCmd.AddCommand(domainRecordDeleteCmd)
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
		domainID, err := resolveDomain(args[0])
		if err != nil {
			return err
		}

		recordIDs, err := convertStringsToIDs(args[1:])
		if err != nil {
			return err
		}

		records, err := getDomainRecords(domainID, recordIDs)
		if err != nil {
			return err
		}

		outputFormatTable(records, []string{"ID", "TYPE", "NAME", "CONTENT"}, []string{"ID", "Type", "Name", "Content"})
		return nil
	},
}

func getDomainRecords(domainID l27.IntID, ids []l27.IntID) ([]l27.DomainRecord, error) {
	c := Level27Client
	if len(ids) == 0 {
		return c.DomainRecords(domainID, recordGetType, optGetParameters)
	} else {
		domains := make([]l27.DomainRecord, len(ids))
		for idx, id := range ids {
			var err error
			domains[idx], err = c.DomainRecord(domainID, id)
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
var domainRecordCreatePriority int32

var domainRecordCreateCmd = &cobra.Command{
	Use:   "create [domain]",
	Short: "Create a new record for a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := resolveDomain(args[0])
		if err != nil {
			return err
		}

		record, err := Level27Client.DomainRecordCreate(id, l27.DomainRecordRequest{
			Name:     domainRecordCreateName,
			Type:     domainRecordCreateType,
			Priority: domainRecordCreatePriority,
			Content:  domainRecordCreateContent,
		})

		if err != nil {
			return err
		}

		outputFormatTemplate(record, "templates/entities/domainRecord/create.tmpl")

		return nil
	},
}

// DELETE DOMAIN/RECORD
var domainRecordDeleteCmd = &cobra.Command{
	Use:   "delete [domain] [record]",
	Short: "Delete a record for a domain",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		//check for valid domain id
		domainID, err := resolveDomain(args[0])
		if err != nil {
			return err
		}

		recordID, err := checkSingleIntID(args[1], "record")
		if err != nil {
			return err
		}

		err = Level27Client.DomainRecordDelete(domainID, recordID)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/domainRecord/delete.tmpl")

		return nil
	},
}

var domainRecordUpdateName string
var domainRecordUpdateContent string
var domainRecordUpdatePriority int32

var domainRecordUpdateCmd = &cobra.Command{
	Use:   "update [domain] [record]",
	Short: "Update a record for a domain",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		domainID, err := resolveDomain(args[0])
		if err != nil {
			return err
		}

		recordID, err := checkSingleIntID(args[1], "record")
		if err != nil {
			return err
		}

		// Merge data with existing so we don't bulldoze anything.
		data, err := Level27Client.DomainRecord(domainID, recordID)
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

		err = Level27Client.DomainRecordUpdate(domainID, recordID, request)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/domainRecord/update.tmpl")
		return nil
	},
}
