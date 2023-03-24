package cmd

import (
	"fmt"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

func init() {
	// APP COMPONENT URL
	appComponentCmd.AddCommand(appComponentUrlCmd)

	// APP COMPONENT URL GET
	appComponentUrlCmd.AddCommand(appComponentUrlGetCmd)
	addCommonGetFlags(appComponentUrlGetCmd)

	// APP COMPONENT URL CREATE
	appComponentUrlCmd.AddCommand(appComponentUrlCreateCmd)
	addWaitFlag(appComponentUrlCreateCmd)
	appComponentUrlCreateCmd.Flags().BoolVar(&appComponentUrlCreateAuthentication, "authentication", false, "Require HTTP Basic authentication on the URL")
	appComponentUrlCreateCmd.Flags().StringVarP(&appComponentUrlCreateContent, "content", "c", "", "Content for the new URL")
	appComponentUrlCreateCmd.Flags().BoolVar(&appComponentUrlCreateSslForce, "force-ssl", false, "Force usage of SSL on the URL")
	appComponentUrlCreateCmd.Flags().Int32Var(&appComponentUrlCreateSslCertificate, "ssl-certificate", 0, "SSL certificate to use.")
	appComponentUrlCreateCmd.Flags().BoolVar(&appComponentUrlCreateHandleDns, "handle-dns", false, "Automatically create DNS records")
	appComponentUrlCreateCmd.Flags().BoolVar(&appComponentUrlCreateAutoSslCertificate, "auto-ssl-certificate", false, "Automatically create SSL certificate with Let's Encrypt")
	appComponentUrlCreateCmd.Flags().BoolVar(&appComponentUrlCreateCaching, "caching", true, "Whether to enable caching on the URL, if the component has a linked Varnish component")

	// APP COMPONENT URL DELETE
	appComponentUrlCmd.AddCommand(appComponentUrlDeleteCmd)
	addWaitFlag(appComponentUrlDeleteCmd)
	appComponentUrlDeleteCmd.Flags().BoolVar(&appComponentUrlDeleteForce, "force", false, "Do not ask for confirmation to delete the URL")
}

// Resolve the ID of an app component URL based on user-provided name or ID.
func resolveAppComponentUrl(appID l27.IntID, appComponentID l27.IntID, arg string) (l27.IntID, error) {
	// if arg already int, this is the ID
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.AppComponentUrlLookup(appID, appComponentID, arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"url",
		func(url l27.AppComponentUrlShort) string { return fmt.Sprintf("%s (%d)", url.Content, url.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, err
}

// ------- APP COMPONENT URLs

// APP COMPONENT URL
var appComponentUrlCmd = &cobra.Command{
	Use:     "url",
	Aliases: []string{"urls"},
}

// APP COMPONENT URL GET
var appComponentUrlGetCmd = &cobra.Command{
	Use: "get",

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

		results, err := resolveGets(
			args[2:],
			func(name string) ([]l27.AppComponentUrlShort, error) {
				return Level27Client.AppComponentUrlLookup(appID, componentID, name)
			},
			func(i l27.IntID) (l27.AppComponentUrlShort, error) {
				res, err := Level27Client.AppComponentUrlGetSingle(appID, componentID, i)
				if err != nil {
					return l27.AppComponentUrlShort{}, err
				}
				return res.ToShort(), nil
			},
			func(cgp l27.CommonGetParams) ([]l27.AppComponentUrlShort, error) {
				return Level27Client.AppComponentUrlGetList(appID, componentID, cgp)
			},
		)

		if err != nil {
			return err
		}

		outputFormatTable(
			results,
			[]string{"ID", "CONTENT", "STATUS", "TYPE", "SSL CERT", "FORCE SSL", "HANDLE DNS", "AUTHENTICATE", "CACHING"},
			[]string{"ID", "Content", "Status", "Type", "SslCertificate.Name", "SslForce", "HandleDNS", "Authentication", "Caching"})

		return nil
	},
}

// APP COMPONENT URL CREATE
var appComponentUrlCreateAuthentication bool
var appComponentUrlCreateContent string
var appComponentUrlCreateSslForce bool
var appComponentUrlCreateSslCertificate l27.IntID
var appComponentUrlCreateHandleDns bool
var appComponentUrlCreateAutoSslCertificate bool
var appComponentUrlCreateCaching bool
var appComponentUrlCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an url for an appcomponent.",

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

		var cert *l27.IntID
		if appComponentUrlCreateSslCertificate == 0 {
			cert = nil
		} else {
			cert = &appComponentUrlCreateSslCertificate
		}

		create := l27.AppComponentUrlCreate{
			Authentication:     appComponentUrlCreateAuthentication,
			Content:            appComponentUrlCreateContent,
			SslForce:           appComponentUrlCreateSslForce,
			SslCertificate:     cert,
			HandleDns:          appComponentUrlCreateHandleDns,
			AutoSslCertificate: appComponentUrlCreateAutoSslCertificate,
			Caching:            appComponentUrlCreateCaching,
		}

		url, err := Level27Client.AppComponentUrlCreate(appID, componentID, create)
		if err != nil {
			return err
		}

		if optWait {
			url, err = waitForStatus(
				func() (l27.AppComponentUrl, error) {
					return Level27Client.AppComponentUrlGetSingle(appID, componentID, url.ID)
				},
				func(s l27.AppComponentUrl) string { return s.Status },
				"ok",
				[]string{"creating"},
			)

			if err != nil {
				return fmt.Errorf("waiting on URL status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(url, "templates/entities/appComponentUrl/create.tmpl")
		return nil
	},
}

// APP COMPONENT URL DELETE
var appComponentUrlDeleteForce bool
var appComponentUrlDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an url from an appcomponent.",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		componentID, err := resolveAppComponent(appID, args[1])
		if err != nil {
			return err
		}

		urlID, err := resolveAppComponentUrl(appID, componentID, args[2])
		if err != nil {
			return err
		}

		if !appComponentUrlDeleteForce {
			url, err := Level27Client.AppComponentUrlGetSingle(appID, componentID, urlID)
			if err != nil {
				return err
			}

			msg := fmt.Sprintf(
				"Delete URL %s (%d) on app comp %s (%d)?",
				url.Content, url.ID,
				url.Appcomponent.Name, url.Appcomponent.ID)

			if !confirmPrompt(msg) {
				return nil
			}
		}

		err = Level27Client.AppComponentUrlDelete(appID, componentID, urlID)
		if err != nil {
			return err
		}

		if optWait {
			err = waitForDelete(
				func() (l27.AppComponentUrl, error) {
					return Level27Client.AppComponentUrlGetSingle(appID, componentID, urlID)
				},
				func(a l27.AppComponentUrl) string { return a.Status },
				[]string{"to_delete"},
			)

			if err != nil {
				return fmt.Errorf("waiting on app status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(nil, "templates/entities/appComponentUrl/delete.tmpl")
		return nil
	},
}
