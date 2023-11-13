package cmd

import (
	"fmt"
	"strings"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

func init() {
	// LVL PACKAGE
	RootCmd.AddCommand(cmdCustomPackages)

	// LVL PACKAGE BASE
	cmdCustomPackages.AddCommand(cmdCustomPackagesBase)

	// LVL PACKAGE BASE GET
	cmdCustomPackagesBase.AddCommand(cmdCustomPackagesBaseGet)

	// LVL PACKAGE CREATE
	cmdCustomPackages.AddCommand(customPackagesCreateCmd)
	customPackagesCreateCmd.Flags().StringVarP(&customPackagesCreateBase, "base", "b", "", "The base package to use. Use 'lvl package base get' to show base packages.")
	customPackagesCreateCmd.Flags().StringVarP(&customPackagesCreateName, "name", "n", "", "The name that the created package will have.")
	customPackagesCreateCmd.Flags().StringVar(&customPackagesCreateOrganisation, "organisation", "", "Organisation that the package will be created on. Defaults to your current organisation.")
	_ = customPackagesCreateCmd.MarkFlagRequired("base")
	_ = customPackagesCreateCmd.MarkFlagRequired("name")

	// LVL PACKAGE
	cmdCustomPackages.AddCommand(customPackagesGetCmd)
	addCommonGetFlags(customPackagesGetCmd)

	// LVL PACKAGE DESCRIBE
	cmdCustomPackages.AddCommand(customPackagesDescribeCmd)

	// LVL PACKAGE INSTANTIATE
	cmdCustomPackages.AddCommand(customPackageInstantiateCmd)
	customPackageInstantiateCmd.Flags().StringVar(&customPackageInstantiateOrganisation, "organisation", "", "")
	customPackageInstantiateCmd.Flags().StringArrayVarP(&customPackageInstantiateParams, "param", "p", nil, "")
	customPackageInstantiateCmd.Flags().StringVarP(&customPackageInstantiateName, "name", "n", "", "Name of the created app")
	customPackageInstantiateCmd.MarkFlagRequired(customPackageInstantiateName)
	addWaitFlag(customPackageInstantiateCmd)

	// LVL PACKAGE DELETE
	cmdCustomPackages.AddCommand(customPackagesDeleteCmd)
	addDeleteConfirmFlag(customPackagesDeleteCmd)

	// LVL PACKAGE TEMPLATE
	cmdCustomPackages.AddCommand(customPackageTemplateCmd)

	// LVL PACKAGE TEMPLATE ADD
	customPackageTemplateCmd.AddCommand(customPackageTemplateAddCmd)
	customPackageTemplateAddCmd.Flags().StringVarP(&customPackageTemplateAddTemplate, "template", "t", "", "")
	customPackageTemplateAddCmd.Flags().StringVarP(&customPackageTemplateAddGroup, "group", "g", "", "")
	_ = customPackagesCreateCmd.MarkFlagRequired("template")
	_ = customPackagesCreateCmd.MarkFlagRequired("group")

	// LVL PACKAGE TEMPLATE REMOVE
	customPackageTemplateCmd.AddCommand(customPackageTemplateRemoveCmd)

	// LVL PACKAGE TEMPLATE TYPE
	customPackageTemplateCmd.AddCommand(customPackageTemplateTypeCmd)

	// LVL PACKAGE TEMPLATE TYPE GET
	customPackageTemplateTypeCmd.AddCommand(customPackageTemplateTypeGetCmd)

	// LVL PACKAGE TEMPLATE TYPE DESCRIBE
	customPackageTemplateTypeCmd.AddCommand(customPackageTemplateTypeDescribeCmd)
}

// Resolve the ID of an app based on user-provided name or ID.
func resolveCustomPackage(arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.CustomPackageLookup(arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"custom package",
		func(pack l27.CustomPackageShort) string { return fmt.Sprintf("%s (%d)", pack.Name, pack.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, nil
}

var cmdCustomPackages = &cobra.Command{
	Use:   "package",
	Short: "Manage Agency Hosting packages",
	Long: `Manage Agency Hosting packages.
Packages are, fundamentally, a series of templates that get executed to create the final set of products.`,
	Aliases: []string{"packages", "pkg", "pkgs", "custompackage", "custompackages"},
}

var cmdCustomPackagesBase = &cobra.Command{
	Use:   "base",
	Short: "Commands for showing base packages",
}

var cmdCustomPackagesBaseGet = &cobra.Command{
	Use:     "get",
	Short:   "List base packages that Agency packages can be created from",
	Example: "  lvl package base get",

	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		packages, err := Level27Client.CustPackagesGetList()
		if err != nil {
			return err
		}

		outputFormatTable(
			packages,
			[]string{"NAME", "DISPLAY NAME"},
			[]string{"Name", "DisplayName"})

		return nil
	},
}

var customPackagesCreateBase string
var customPackagesCreateName string
var customPackagesCreateOrganisation string
var customPackagesCreateCmd = &cobra.Command{
	Use:     "create -b <base package> -n <name>",
	Short:   "Create a new package",
	Example: `  lvl package create -b php_mysql_generic_walk -n walk-php`,

	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		orgID, err := resolveOrgOrUserOrg(customPackagesCreateOrganisation)
		if err != nil {
			return err
		}

		create := l27.CustomPackageCreate{
			Name:              customPackagesCreateName,
			CustomPackageName: customPackagesCreateBase,
			Organisation:      orgID,
		}

		pack, err := Level27Client.CustomPackageCreate(&create)
		if err != nil {
			return err
		}

		outputFormatTemplate(pack, "templates/entities/customPackages/create.tmpl")

		return nil
	},
}

var customPackagesDeleteCmd = &cobra.Command{
	Use:     "delete <package name or ID>",
	Short:   "Delete a package",
	Long:    "Delete a package. Deleting a package does not affect created apps.",
	Example: "  lvl package delete walk-php",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageID, err := resolveCustomPackage(args[0])
		if err != nil {
			return err
		}

		if !optDeleteConfirmed {
			pack, err := Level27Client.CustomPackageGetSingle(packageID)
			if err != nil {
				return err
			}

			if !confirmPrompt(fmt.Sprintf("Delete package %s (%d)?", pack.Name, pack.ID)) {
				return nil
			}
		}

		err = Level27Client.CustomPackageDelete(packageID)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/customPackages/delete.tmpl")

		return nil
	},
}

var customPackagesGetCmd = &cobra.Command{
	Use:   "get [package name or ID [package...]]",
	Short: "List packages or get info",
	Example: `List packages:
  lvl package get
List packages with search:
  lvl package get -f walk
Get info about a specific package:
  lvl package get walk-php
`,

	RunE: func(cmd *cobra.Command, args []string) error {
		packages, err := resolveGets[l27.CustomPackageShort](
			args,
			func(s string) ([]l27.CustomPackageShort, error) {
				return Level27Client.CustomPackageLookup(s)
			},
			func(i int32) (l27.CustomPackageShort, error) {
				value, err := Level27Client.CustomPackageGetSingle(i)
				if err != nil {
					return l27.CustomPackageShort{}, err
				}

				return value.ToShort(), nil
			},
			func(cgp l27.CommonGetParams) ([]l27.CustomPackageShort, error) {
				return Level27Client.CustomPackageGetList(cgp)
			},
		)

		if err != nil {
			return err
		}

		outputFormatTable(
			packages,
			[]string{"ID", "NAME", "TYPE"},
			[]string{"ID", "Name", "Type"})

		return nil
	},
}

var customPackagesDescribeCmd = &cobra.Command{
	Use:     "describe <package name or ID>",
	Short:   "Show detailed information about a package",
	Example: "  lvl package describe walk-php",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageName := args[0]
		id, err := resolveCustomPackage(packageName)
		if err != nil {
			return err
		}

		pack, err := Level27Client.CustomPackageGetSingle(id)
		if err != nil {
			return err
		}

		outputFormatTemplate(pack, "templates/entities/customPackages/describe.tmpl")

		return nil
	},
}

var customPackageInstantiateOrganisation string
var customPackageInstantiateParams []string
var customPackageInstantiateName string
var customPackageInstantiateCmd = &cobra.Command{
	Use:   "instantiate <package name or ID> -n <name>",
	Short: "Create an app from a package",
	Long: `Create an app from a package.
This only starts the creation via a task. If you wait for the task to complete with --wait,
the resulting entities will be listed.`,

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageName := args[0]
		id, err := resolveCustomPackage(packageName)
		if err != nil {
			return err
		}

		orgID, err := resolveOrgOrUserOrg(customPackageInstantiateOrganisation)
		if err != nil {
			return err
		}

		params := map[string]l27.ParameterValue{}
		params["name"] = customPackageInstantiateName

		for _, param := range customPackageInstantiateParams {
			split := strings.SplitN(param, "=", 2)
			if len(split) != 2 {
				return fmt.Errorf("expected key=value pair to --param: %s", param)
			}

			params[split[0]], err = readArgFileSupported(split[1])
			if err != nil {
				return err
			}
		}

		request := l27.CustomPackageRootTaskRequest{
			Organisation: orgID,
			Params:       params,
		}

		task, err := Level27Client.CustomPackageRootTask(id, &request)
		if err != nil {
			return err
		}

		outputFormatTemplate(task, "templates/entities/customPackages/instantiate.tmpl")

		if optWait {
			task, err = waitForStatus[l27.RootTask](
				func() (l27.RootTask, error) {
					return Level27Client.RootTaskGetSingle(task.Id)
				},
				func(rt l27.RootTask) string { return rt.Status },
				"done",
				[]string{"to_do", "busy"},
			)

			if err != nil {
				return fmt.Errorf("waiting on task failed: %s", err.Error())
			}

			outputFormatTemplate(task, "templates/entities/customPackages/instantiate_full.tmpl")
		}

		return nil
	},
}

var customPackageTemplateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage templates on custom packages",
	// TODO: template add can't even add parameters so it's pretty useless right now
	// Hiding as a result.
	Hidden: true,
}

var customPackageTemplateTypeCmd = &cobra.Command{
	Use:   "type",
	Short: "Show information for available templates",
}

var customPackageTemplateTypeGetCmd = &cobra.Command{
	Use:   "get",
	Short: "List templates available for adding to an Agency package",

	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		templates, err := Level27Client.TemplatesGetList(true)
		if err != nil {
			return err
		}

		outputFormatTable(templates, []string{"NAME", "DISPLAY NAME"}, []string{"Name", "DisplayName"})

		return nil
	},
}

var customPackageTemplateTypeDescribeCmd = &cobra.Command{
	Use:   "describe <template name>",
	Short: "Show detailed information for a package template",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		template, err := Level27Client.TemplatesGetSingle(args[0], l27.TemplateGetSingleRequest{CustomPackagePossible: true})
		if err != nil {
			return err
		}

		outputFormatTemplate(template, "templates/entities/template/describe.tmpl")

		return nil
	},
}

var customPackageTemplateAddTemplate string
var customPackageTemplateAddGroup string
var customPackageTemplateAddCmd = &cobra.Command{
	Use:   "add <package> -t <template>",
	Short: "Add a template to a package",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		customPackageID, err := resolveCustomPackage(args[0])
		if err != nil {
			return err
		}

		create := l27.CustomPackageTemplateCreate{
			Template:      customPackageTemplateAddTemplate,
			LimitGroup:    customPackageTemplateAddGroup,
			CustomPackage: customPackageID,
		}

		template, err := Level27Client.CustomPackageTemplateCreate(customPackageID, &create)
		if err != nil {
			return err
		}

		outputFormatTemplate(template, "templates/entities/customPackageTemplate/create.tmpl")

		return nil
	},
}

var customPackageTemplateRemoveCmd = &cobra.Command{
	Use: "remove",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		customPackageID, err := resolveCustomPackage(args[0])
		if err != nil {
			return err
		}

		templateID, err := checkSingleIntID(args[1], "template")
		if err != nil {
			return err
		}

		err = Level27Client.CustomPackageTemplateRemove(customPackageID, templateID)
		if err != nil {
			return err
		}

		return nil
	},
}
