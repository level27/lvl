package cmd

import (
	"fmt"
	"strings"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(taskCmd)

	taskCmd.AddCommand(taskCreateCmd)
	taskCreateCmd.Flags().StringVar(&taskCreateOrganisation, "organisation", "", "")
	taskCreateCmd.Flags().StringVarP(&taskCreateTemplate, "template", "t", "", "")
	taskCreateCmd.Flags().StringVarP(&taskCreatePackage, "package", "P", "", "")
	taskCreateCmd.Flags().StringArrayVarP(&taskCreateParameters, "param", "p", nil, "")

	taskCmd.AddCommand(taskDescribeCmd)
}

var taskCmd = &cobra.Command{
	Use:     "task",
	Aliases: []string{"tasks"},
	Short:   "Commands for managing tasks",
}

var taskCreatePackage string
var taskCreateTemplate string
var taskCreateOrganisation string
var taskCreateParameters []string

var taskCreateCmd = &cobra.Command{
	Use:  "create",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		org, err := resolveOrgOrUserOrg(appComponentAttachmentUploadOrganisation)
		if err != nil {
			return fmt.Errorf("couldn't resolve organisation: %v", err)
		}

		var template, pkg *string
		if taskCreateTemplate != "" {
			if taskCreatePackage != "" {
				return fmt.Errorf("cannot specify both a template and a package")
			}

			template = &taskCreateTemplate
		} else if taskCreatePackage != "" {
			pkg = &taskCreatePackage
		} else {
			return fmt.Errorf("must specify either a template or a package")
		}

		params := map[string]l27.ParameterValue{}

		for _, param := range taskCreateParameters {
			split := strings.SplitN(param, "=", 2)
			if len(split) != 2 {
				return fmt.Errorf("expected key=value pair to --param: %s", param)
			}

			params[split[0]], err = readArgFileSupported(split[1])
			if err != nil {
				return err
			}
		}

		request := l27.RootTaskCreate{
			Organisation: org,
			Template:     template,
			Package:      pkg,
			Parameters:   params,
		}

		task, err := Level27Client.RootTaskCreate(request)
		if err != nil {
			return fmt.Errorf("failed to create task: %v", err)
		}

		outputFormatTemplate(task, "templates/entities/task/create.tmpl")

		return nil
	},
}

var taskDescribeCmd = &cobra.Command{
	Use:  "describe <root task ID>",
	Args: cobra.ExactArgs(1),
	// Temporary thing I made to test some templates.
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := l27.ParseID(args[0])
		if err != nil {
			return err
		}

		task, err := Level27Client.RootTaskGetSingle(id)
		if err != nil {
			return err
		}

		outputFormatTemplate(task, "templates/entities/customPackages/instantiate_full.tmpl")

		return nil
	},
}
