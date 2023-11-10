package cmd

import (
	"fmt"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

func init() {
	appComponentCmd.AddCommand(appComponentDomainCmd)

	appComponentDomainCmd.AddCommand(appComponentDomainLinkCmd)
	addWaitFlag(appComponentDomainLinkCmd)
	appComponentDomainLinkCmd.Flags().StringVarP(&optAppComponentDomainLinkDomain, "domain", "d", "", "Domain to link to")
	appComponentDomainLinkCmd.Flags().BoolVar(&optAppComponentDomainLinkHandleDNS, "handle-dns", true, "Automatically manage relevant DNS records")
	appComponentDomainLinkCmd.Flags().BoolVar(&optAppComponentDomainLinkDKIM, "dkim", true, "Create DKIM keys for mail components")
	appComponentDomainLinkCmd.MarkFlagRequired("domain")

	appComponentDomainCmd.AddCommand(appComponentDomainUnlinkCmd)
	addWaitFlag(appComponentDomainUnlinkCmd)

	appComponentDomainCmd.AddCommand(appComponentDomainGetCmd)
	addCommonGetFlags(appComponentDomainGetCmd)
}

// APP COMPONENT DOMAIN
var appComponentDomainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Commands for managing linked domains",
	Long:  "Components like mail allow linking of domains to automatically manage DNS and similar.",
}

var appComponentDomainGetCmd = &cobra.Command{
	Use:   "get <app> <component> [domain [domain...]]",
	Short: "List domains linked to an app component",
	Example: `List linked domains:
  lvl app component domain get my-app mail
Get info about a single linked domain:
  lvl app component domain get my-app mail example.com`,

	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		componentID, err := resolveAppComponent(appID, args[1])
		if err != nil {
			return err
		}

		domains, err := resolveGets(
			args[2:],
			func(s string) ([]l27.AppComponentDomainShort, error) {
				return Level27Client.AppComponentDomainLookup(appID, componentID, s)
			},
			func(id l27.IntID) (l27.AppComponentDomainShort, error) {
				value, err := Level27Client.AppComponentDomainGetSingle(appID, componentID, id)
				if err != nil {
					return l27.AppComponentDomainShort{}, err
				}

				return value.ToShort(), nil
			},
			func(cgp l27.CommonGetParams) ([]l27.AppComponentDomainShort, error) {
				return Level27Client.AppComponentDomainGetList(appID, componentID, cgp)
			},
		)

		if err != nil {
			return err
		}

		outputFormatTableFuncs(
			domains,
			[]string{"ID", "NAME", "STATUS", "HANDLE DNS", "DKIM"},
			[]interface{}{
				"Domain.ID",
				func(dom l27.AppComponentDomainShort) string {
					return fmt.Sprintf("%s.%s", dom.Domain.Name, dom.Domain.Domaintype.Extension)
				},
				"Status",
				"HandleDNS",
				"DKIM",
			})

		return nil
	},
}

var optAppComponentDomainLinkDomain string
var optAppComponentDomainLinkHandleDNS bool
var optAppComponentDomainLinkDKIM bool
var appComponentDomainLinkCmd = &cobra.Command{
	Use:   "link <app> <component> -d <domain>",
	Short: "Link a domain to a component",
	Example: `Link a domain:
  lvl app component domain link my-app mail -d example.com
Link a domain, without handling DNS automatically:
  lvl app component domain link my-app mail -d example.com --handle-dns=false
`,

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		componentID, err := resolveAppComponent(appID, args[1])
		if err != nil {
			return err
		}

		domainID, err := resolveDomain(optAppComponentDomainLinkDomain)
		if err != nil {
			return err
		}

		create := l27.AppComponentDomainCreate{
			Domain:    domainID,
			HandleDNS: optAppComponentDomainLinkHandleDNS,
			DKIM:      optAppComponentDomainLinkDKIM,
		}

		domain, err := Level27Client.AppComponentDomainCreate(appID, componentID, create)
		if err != nil {
			return fmt.Errorf("linking domain failed: %s", err.Error())
		}

		if optWait {
			domain, err = waitForStatus(
				func() (l27.AppComponentDomain, error) {
					return Level27Client.AppComponentDomainGetSingle(appID, componentID, domain.ID)
				},
				func(s l27.AppComponentDomain) string { return s.Status },
				"ok",
				[]string{"to_create", "creating"},
			)

			if err != nil {
				return fmt.Errorf("waiting on domain status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(domain, "templates/entities/appComponentDomain/create.tmpl")

		return nil
	},
}

var appComponentDomainUnlinkCmd = &cobra.Command{
	Use:   "unlink <app> <component> <domain>",
	Short: "Unlink a domain from a component",
	Example: `Unlink a domain:
  lvl app component domain unlink my-app mail example.com`,

	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		componentID, err := resolveAppComponent(appID, args[1])
		if err != nil {
			return err
		}

		domainID, err := resolveDomain(args[2])
		if err != nil {
			return err
		}

		err = Level27Client.AppComponentDomainDelete(appID, componentID, domainID)
		if err != nil {
			return fmt.Errorf("unlinking domain failed: %s", err.Error())
		}

		if optWait {
			err = waitForDelete(
				func() (l27.AppComponentDomain, error) {
					return Level27Client.AppComponentDomainGetSingle(appID, componentID, domainID)
				},
				func(s l27.AppComponentDomain) string { return s.Status },
				[]string{"deleting"},
			)

			if err != nil {
				return fmt.Errorf("waiting on domain status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(nil, "templates/entities/appComponentDomain/delete.tmpl")

		return nil
	},
}
