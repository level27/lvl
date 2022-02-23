package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var jobCmd = &cobra.Command{
	Use: "job",
	Short: "Commands related to viewing and managing jobs",
}

var jobDescribeCmd = &cobra.Command{
	Use: "describe",
	Short: "Get complete overview of a job",
	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		jobId, err := convertStringToId(args[0])
		if err != nil {
			log.Fatalln("Invalid job ID!")
		}

		job := Level27Client.JobHistoryRootGet(jobId)
		outputFormatTemplate(job, "templates/job.tmpl")
	},
}

func init() {
	RootCmd.AddCommand(jobCmd)

	jobCmd.AddCommand(jobDescribeCmd)
}