package cmd

import (
	"fmt"
	"strings"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

func init() {
	// LVL CUSTOMPACKAGES
	RootCmd.AddCommand(cmdCustomPackages)

	// LVL CUSTOMPACKAGES BASE
	cmdCustomPackages.AddCommand(cmdCustomPackagesBase)

	// LVL CUSTOMPACKAGES BASE GET
	cmdCustomPackagesBase.AddCommand(cmdCustomPackagesBaseGet)

	// LVL CUSTOMPACKAGES CREATE
	cmdCustomPackages.AddCommand(customPackagesCreateCmd)
	customPackagesCreateCmd.Flags().StringVarP(&customPackagesCreateBase, "base", "b", "", "")
	customPackagesCreateCmd.Flags().StringVarP(&customPackagesCreateName, "name", "n", "", "")
	customPackagesCreateCmd.Flags().StringVar(&customPackagesCreateOrganisation, "organisation", "", "")
	_ = customPackagesCreateCmd.MarkFlagRequired("base")
	_ = customPackagesCreateCmd.MarkFlagRequired("name")

	// LVL CUSTOMPACKAGES
	cmdCustomPackages.AddCommand(customPackagesGetCmd)
	addCommonGetFlags(customPackagesGetCmd)

	// LVL CUSTOMPACKAGES DESCRIBE
	cmdCustomPackages.AddCommand(customPackagesDescribeCmd)

	// LVL CUSTOMPACKAGES INSTANTIATE
	cmdCustomPackages.AddCommand(customPackageInstantiateCmd)
	customPackageInstantiateCmd.Flags().StringVar(&customPackageInstantiateOrganisation, "organisation", "", "")
	customPackageInstantiateCmd.Flags().StringArrayVarP(&customPackageInstantiateParams, "param", "p", nil, "")

	// LVL CUSTOMPACKAGES DELETE
	cmdCustomPackages.AddCommand(customPackagesDeleteCmd)
	addDeleteConfirmFlag(customPackagesDeleteCmd)

	// LVL CUSTOMPACKAGES TEMPLATE
	cmdCustomPackages.AddCommand(customPackageTemplateCmd)

	// LVL CUSTOMPACKAGES TEMPLATE ADD
	customPackageTemplateCmd.AddCommand(customPackageTemplateAddCmd)
	customPackageTemplateAddCmd.Flags().StringVarP(&customPackageTemplateAddTemplate, "template", "t", "", "")
	customPackageTemplateAddCmd.Flags().StringVarP(&customPackageTemplateAddGroup, "group", "g", "", "")
	_ = customPackagesCreateCmd.MarkFlagRequired("template")
	_ = customPackagesCreateCmd.MarkFlagRequired("group")

	// LVL CUSTOMPACKAGES TEMPLATE REMOVE
	customPackageTemplateCmd.AddCommand(customPackageTemplateRemoveCmd)

	// LVL CUSTOMPACKAGES TEMPLATE TYPE
	customPackageTemplateCmd.AddCommand(customPackageTemplateTypeCmd)

	// LVL CUSTOMPACKAGES TEMPLATE TYPE GET
	customPackageTemplateTypeCmd.AddCommand(customPackageTemplateTypeGetCmd)

	// LVL CUSTOMPACKAGES TEMPLATE TYPE DESCRIBE
	customPackageTemplateTypeCmd.AddCommand(customPackageTemplateTypeDescribeCmd)
}

// Resolve the ID of an app based on user-provided name or ID.
func resolveCustomerPackage(arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.CustomerPackageLookup(arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"custom package",
		func(pack l27.CustomerPackageShort) string { return fmt.Sprintf("%s (%d)", pack.Name, pack.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, nil
}

var cmdCustomPackages = &cobra.Command{
	Use:     "custompackages",
	Aliases: []string{"custpkgs"},
	Hidden:  true,
}

var cmdCustomPackagesBase = &cobra.Command{
	Use:   "base",
	Short: "Commands for showing base custom packages",
}

var cmdCustomPackagesBaseGet = &cobra.Command{
	Use:   "get",
	Short: "List base custom packages",

	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		packages, err := Level27Client.CustomPackagesGetList()
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
	Use: "create -b <base package> -n <name>",

	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		orgID, err := resolveOrgOrUserOrg(customPackagesCreateOrganisation)
		if err != nil {
			return err
		}

		create := l27.CustomerPackageCreate{
			Name:              customPackagesCreateName,
			CustomPackageName: customPackagesCreateBase,
			Organisation:      orgID,
		}

		pack, err := Level27Client.CustomerPackageCreate(&create)
		if err != nil {
			return err
		}

		outputFormatTemplate(pack, "templates/entities/customPackages/create.tmpl")

		return nil
	},
}

var customPackagesDeleteCmd = &cobra.Command{
	Use: "delete <custom package>",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageID, err := resolveCustomerPackage(args[0])
		if err != nil {
			return err
		}

		if !optDeleteConfirmed {
			pack, err := Level27Client.CustomerPackageGetSingle(packageID)
			if err != nil {
				return err
			}

			if !confirmPrompt(fmt.Sprintf("Delete custom package %s (%d)?", pack.Name, pack.ID)) {
				return nil
			}
		}

		err = Level27Client.CustomerPackageDelete(packageID)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/customPackages/delete.tmpl")

		return nil
	},
}

var customPackagesGetCmd = &cobra.Command{
	Use: "get",

	RunE: func(cmd *cobra.Command, args []string) error {
		packages, err := Level27Client.CustomerPackageGetList(optGetParameters)
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
	Use: "describe <custom package>",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageName := args[0]
		id, err := resolveCustomerPackage(packageName)
		if err != nil {
			return err
		}

		pack, err := Level27Client.CustomerPackageGetSingle(id)
		if err != nil {
			return err
		}

		outputFormatTemplate(pack, "templates/entities/customPackages/describe.tmpl")

		return nil
	},
}

var customPackageInstantiateOrganisation string
var customPackageInstantiateParams []string
var customPackageInstantiateCmd = &cobra.Command{
	Use: "instantiate <custom package>",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packageName := args[0]
		id, err := resolveCustomerPackage(packageName)
		if err != nil {
			return err
		}

		orgID, err := resolveOrgOrUserOrg(customPackageInstantiateOrganisation)
		if err != nil {
			return err
		}

		params := map[string]l27.ParameterValue{}

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

		request := l27.CustomerPackageRootTaskRequest{
			Organisation: orgID,
			Params:       params,
		}

		err = Level27Client.CustomerPackageRootTask(id, &request)
		if err != nil {
			return err
		}

		return nil
	},
}

var customPackageTemplateCmd = &cobra.Command{
	Use: "template",
}

var customPackageTemplateTypeCmd = &cobra.Command{
	Use: "type",
}

var customPackageTemplateTypeGetCmd = &cobra.Command{
	Use: "get",

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
	Use: "describe",

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
	Use: "add",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		customPackageID, err := resolveCustomerPackage(args[0])
		if err != nil {
			return err
		}

		create := l27.CustomerPackageTemplateCreate{
			Template:      customPackageTemplateAddTemplate,
			LimitGroup:    customPackageTemplateAddGroup,
			CustomPackage: customPackageID,
		}

		template, err := Level27Client.CustomerPackageTemplateCreate(customPackageID, &create)
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
		customPackageID, err := resolveCustomerPackage(args[0])
		if err != nil {
			return err
		}

		templateID, err := checkSingleIntID(args[1], "template")
		if err != nil {
			return err
		}

		err = Level27Client.CustomerPackageTemplateRemove(customPackageID, templateID)
		if err != nil {
			return err
		}

		return nil
	},
}
