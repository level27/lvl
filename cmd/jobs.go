package cmd

import (
	"log"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(jobCmd)

	jobCmd.AddCommand(jobDescribeCmd)
}

var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Commands related to viewing and managing jobs",
}

var jobDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Get complete overview of a job",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		jobId, err := convertStringToId(args[0])
		if err != nil {
			log.Fatalln("Invalid job ID!")
		}

		job := Level27Client.JobHistoryRootGet(jobId)
		outputFormatTemplate(job, "templates/job.tmpl")
	},
}

// Add common commands for managing entity jobs to a parent command.
// entityType is the type for /jobs/history/{type}/{id} which this function uses.
// resolve is a function that turns an argument in the ID of the entity.
func addJobCmds(parent *cobra.Command, entityType string, resolve func(string) int) {
	jobsCmd := &cobra.Command{
		Use:   "jobs",
		Short: "View job history for this entity",
		Args:  cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			entityID := resolve(args[0])

			//get full history of toplevel jobs
			history := Level27Client.EntityJobHistoryGet(entityType, entityID)

			// filter jobs where status is not 50.
			notCompleted := FindNotcompletedJobs(history)

			// check for every job without status 50. the subjobs who don't have status 50
			for _, RootJob := range notCompleted {
				fullData := Level27Client.JobHistoryRootGet(RootJob.Id)

				for _, subjob := range fullData.Jobs {
					if subjob.Status != 50 {
						notCompleted = append(notCompleted, subjob)
						if len(subjob.Jobs) != 0 {
							notCompleted = append(notCompleted, FindNotcompletedJobs(subjob.Jobs)...)
						}
					}
				}
			}

			outputFormatTable(notCompleted, []string{"ID", "STATUS", "MESSAGE", "DATE"}, []string{"Id", "Status", "Message", "Dt"})
		},
	}

	parent.AddCommand(jobsCmd)
}

func CheckSubJobs(job l27.Job) bool {
	if len(job.Jobs) == 0 {
		return false
	} else {
		return true
	}
}

func FindNotcompletedJobs(jobs []l27.Job) []l27.Job {
	var NotCompleted []l27.Job
	for _, job := range jobs {
		if job.Status != 50 {
			NotCompleted = append(NotCompleted, job)
		}
	}
	return NotCompleted

}
