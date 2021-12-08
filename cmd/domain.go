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
	Short: "Commands for working with domains",
}

var domainGetCmd = &cobra.Command{
	Use: "get",

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

var domainDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Get detailed info about a domain",
	Run: func(cmd *cobra.Command, args []string) {
		Level27Client.DomainDescribe(args)
	},
}

var domainRecordCmd = &cobra.Command{
	Use:   "record",
	Short: "Commands for managing domain records",
}

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


func init() {
	RootCmd.AddCommand(domainCmd)

	domainCmd.AddCommand(domainGetCmd)
	addCommonGetFlags(domainGetCmd)

	domainCmd.AddCommand(domainDescribeCmd)

	domainCmd.AddCommand(domainRecordCmd)
	domainRecordCmd.AddCommand(domainRecordListCmd)

	flags := domainRecordCreateCmd.Flags() 
	flags.StringVarP(&domainRecordCreateType, "type", "t", "", "Type of the domain record")
	flags.StringVarP(&domainRecordCreateName, "name", "n", "", "Name of the domain record")
	flags.StringVarP(&domainRecordCreateContent, "content", "c", "", "Content of the domain record")
	flags.IntVarP(&domainRecordCreatePriority, "priority", "p", 0, "Priority of the domain record")
	domainRecordCreateCmd.MarkFlagRequired("type")
	domainRecordCreateCmd.MarkFlagRequired("content")
	domainRecordCmd.AddCommand(domainRecordCreateCmd)
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