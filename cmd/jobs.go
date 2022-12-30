package cmd

import (
	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(jobCmd)

	jobCmd.AddCommand(jobDescribeCmd)
	jobCmd.AddCommand(jobRetryCmd)
	jobCmd.AddCommand(jobDeleteCmd)
}

var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Commands related to viewing and managing jobs",
}

var jobDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Get complete overview of a job",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		jobID, err := checkSingleIntID(args[0], "job")
		if err != nil {
			return err
		}

		job, err := Level27Client.JobHistoryRootGet(jobID, l27.JobHistoryGetParams{})
		if err != nil {
			return err
		}

		outputFormatTemplate(job, "templates/job.tmpl")
		return err
	},
}

var jobRetryCmd = &cobra.Command{
	Use:     "retry <id>",
	Short:   "(admin only) Retry execution of a job",
	Example: "lvl job retry 12345",
	Args:    cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		jobID, err := checkSingleIntID(args[0], "job")
		if err != nil {
			return err
		}

		err = Level27Client.JobRetry(jobID)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/jobs/retry.tmpl")
		return nil
	},
}

var jobDeleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Short:   "(admin only) Delete a job",
	Example: "lvl job delete 12345",
	Args:    cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		jobID, err := checkSingleIntID(args[0], "job")
		if err != nil {
			return err
		}

		err = Level27Client.JobDelete(jobID)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/jobs/delete.tmpl")
		return nil
	},
}

// Add common commands for managing entity jobs to a parent command.
// entityType is the type for /jobs/history/{type}/{id} which this function uses.
// resolve is a function that turns an argument in the ID of the entity.
func addJobCmds(parent *cobra.Command, entityType string, resolve func(string) (l27.IntID, error)) {
	jobsCmd := &cobra.Command{
		Use:   "jobs",
		Short: "View job history for this entity",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			entityID, err := resolve(args[0])
			if err != nil {
				return err
			}

			//get full history of toplevel jobs
			history, err := Level27Client.EntityJobHistoryGet(entityType, entityID, l27.PageableParams{})
			if err != nil {
				return err
			}

			// filter jobs where status is not 50.
			rootJobs := FindNotcompletedJobsRoot(history)

			shownJobs := []l27.Job{}

			// check for every job without status 50. the subjobs who don't have status 50
			for _, RootJob := range rootJobs {
				fullData, err := Level27Client.JobHistoryRootGet(
					RootJob.ID,
					l27.JobHistoryGetParams{})
				if err != nil {
					return err
				}

				shownJobs = append(shownJobs, fullData)

				for _, subjob := range fullData.Jobs {
					if subjob.Status != 50 {
						shownJobs = append(shownJobs, subjob)
						if len(subjob.Jobs) != 0 {
							shownJobs = append(shownJobs, FindNotcompletedJobs(subjob.Jobs)...)
						}
					}
				}
			}

			outputFormatTable(shownJobs, []string{"ID", "STATUS", "MESSAGE", "DATE"}, []string{"ID", "Status", "Message", "Dt"})
			return nil
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

func FindNotcompletedJobsRoot(jobs []l27.HistoryRootJob) []l27.HistoryRootJob {
	var NotCompleted []l27.HistoryRootJob
	for _, job := range jobs {
		if job.Status != 50 {
			NotCompleted = append(NotCompleted, job)
		}
	}
	return NotCompleted
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
