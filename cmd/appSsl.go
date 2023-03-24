package cmd

import (
	"fmt"

	"github.com/level27/l27-go"
	"github.com/level27/lvl/utils"
	"github.com/spf13/cobra"
)

func init() {
	// APP SSL
	appCmd.AddCommand(appSslCmd)

	// APP SSL GET
	appSslCmd.AddCommand(appSslGetCmd)
	addCommonGetFlags(appSslCmd)

	// APP SSL DESCRIBE
	appSslCmd.AddCommand(appSslDescribeCmd)

	// APP SSL CREATE
	appSslCmd.AddCommand(appSslCreateCmd)
	addWaitFlag(appSslCreateCmd)
	appSslCreateCmd.Flags().StringVarP(&appSslCreateName, "name", "n", "", "Name of this SSL certificate")
	appSslCreateCmd.Flags().StringVarP(&appSslCreateSslType, "type", "t", "", "Type of SSL certificate to use. Options are: letsencrypt, xolphin, own")
	appSslCreateCmd.Flags().StringVar(&appSslCreateAutoSslCertificateUrls, "auto-urls", "", "URL or CSV list of URLs (required for Let's Encrypt)")
	appSslCreateCmd.Flags().StringVar(&appSslCreateSslKey, "ssl-key", "", "SSL key for own certificate. Can be read from a file by specifying @filename.")
	appSslCreateCmd.Flags().StringVar(&appSslCreateSslCabundle, "ssl-cabundle", "", "SSL CA bundle for own certificate. Can be read from a file by specifying @filename.")
	appSslCreateCmd.Flags().StringVar(&appSslCreateSslCrt, "ssl-crt", "", "SSL CRT for own certificate. Can be read from a file by specifying @filename.")
	appSslCreateCmd.Flags().BoolVar(&appSslCreateAutoUrlLink, "auto-link", false, "After creation, automatically link to any URLs without existing certificate")
	appSslCreateCmd.Flags().BoolVar(&appSslCreateSslForce, "ssl-force", false, "Force SSL")
	appSslCreateCmd.MarkFlagRequired("name")
	appSslCreateCmd.MarkFlagRequired("type")

	// APP SSL DELETE
	appSslCmd.AddCommand(appSslDeleteCmd)
	appSslDeleteCmd.Flags().BoolVar(&appSslDeleteForce, "force", false, "Do not ask for confirmation to delete the SSL certificate")

	// APP SSL UPDATE
	appSslCmd.AddCommand(appSslUpdateCmd)
	settingsFileFlag(appSslUpdateCmd)
	settingString(appSslUpdateCmd, updateSettings, "name", "New name for the SSL certificate")

	// APP SSL FIX
	appSslCmd.AddCommand(appSslFixCmd)

	// APP SSL ACTION
	appSslCmd.AddCommand(appSslActionCmd)

	// APP SSL ACTION RETRY
	appSslActionCmd.AddCommand(appSslActionRetryCmd)

	// APP SSL ACTION VALIDATECHALLENGE
	appSslActionCmd.AddCommand(appSslActionValidateChallengeCmd)

	// APP SSL KEY
	appSslCmd.AddCommand(appSslKeyCmd)
}

// Resolve the ID of an SSL certificate based on user-provided name or ID.
func resolveAppSslCertificate(appID l27.IntID, arg string) (l27.IntID, error) {
	// if arg already int, this is the ID
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.AppSslCertificatesLookup(appID, arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"app SSL certificate",
		func(app l27.AppSslCertificate) string { return fmt.Sprintf("%s (%d)", app.Name, app.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, nil
}

// APP SSL CERTIFICATES

// APP SSL
var appSslCmd = &cobra.Command{
	Use:     "ssl",
	Short:   "Commands for managing SSL certificates on apps",
	Example: "lvl app ssl get forum\nlvl app ssl describe forum forum.example.com",

	Aliases: []string{"sslcert"},
}

// APP SSL GET
var appSslGetType string
var appSslGetStatus string
var appSslGetCmd = &cobra.Command{
	Use:     "get [app]",
	Short:   "Get a list of SSL certificates for an app",
	Example: "lvl app ssl get forum\nlvl app ssl get forum -f admin",

	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		certs, err := resolveGets(
			// First arg is app ID.
			args[1:],
			func(name string) ([]l27.AppSslCertificate, error) {
				return Level27Client.AppSslCertificatesLookup(appID, name)
			},
			func(certID l27.IntID) (l27.AppSslCertificate, error) {
				return Level27Client.AppSslCertificatesGetSingle(appID, certID)
			},
			func(get l27.CommonGetParams) ([]l27.AppSslCertificate, error) {
				return Level27Client.AppSslCertificatesGetList(appID, appSslGetType, appSslGetStatus, get)
			},
		)

		if err != nil {
			return err
		}

		outputFormatTableFuncs(
			certs,
			[]string{"ID", "Name", "Type", "Status", "SSL Status", "Expiry Date"},
			[]interface{}{"ID", "Name", "SslType", "Status", "SslStatus", "DtExpires", func(c l27.AppSslCertificate) string { return utils.FormatUnixTime(c.DtExpires) }})

		return nil
	},
}

// APP SSL DESCRIBE
var appSslDescribeCmd = &cobra.Command{
	Use:     "describe [app] [SSL cert]",
	Short:   "Get detailed information of an SSL certificate",
	Example: "lvl app ssl describe forum forum.example.com",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		certID, err := resolveAppSslCertificate(appID, args[1])
		if err != nil {
			return err
		}

		cert, err := Level27Client.AppSslCertificatesGetSingle(appID, certID)
		if err != nil {
			return err
		}

		outputFormatTemplate(cert, "templates/appSslCertificate.tmpl")
		return nil
	},
}

// APP SSL CREATE
var appSslCreateName string
var appSslCreateSslType string
var appSslCreateAutoSslCertificateUrls string
var appSslCreateSslKey string
var appSslCreateSslCrt string
var appSslCreateSslCabundle string
var appSslCreateAutoUrlLink bool
var appSslCreateSslForce bool

var appSslCreateCmd = &cobra.Command{
	Use:     "create [app]",
	Short:   "Create a new SSL certificate on an app",
	Example: "lvl app ssl create forum --name forum.example.com --auto-urls forum.example.com --auto-link --type letsencrypt\nlvl app ssl create forum --name forum.example.com --type own --ssl-cabundle '@cert.ca-bundle' --ssl-key '@key.pem' --ssl-crt '@cert.crt'",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		create := l27.AppSslCertificateCreate{
			Name:                   appSslCreateName,
			SslType:                appSslCreateSslType,
			AutoSslCertificateUrls: appSslCreateAutoSslCertificateUrls,
			SslForce:               appSslCreateSslForce,
			AutoUrlLink:            appSslCreateAutoUrlLink,
		}

		var certificate l27.AppSslCertificate

		switch appSslCreateSslType {
		case "own":
			sslKey, err := readArgFileSupported(appSslCreateSslKey)
			if err != nil {
				return err
			}

			sslCrt, err := readArgFileSupported(appSslCreateSslCrt)
			if err != nil {
				return err
			}

			sslCabundle, err := readArgFileSupported(appSslCreateSslCabundle)
			if err != nil {
				return err
			}

			createOwn := l27.AppSslCertificateCreateOwn{
				AppSslCertificateCreate: create,
				SslKey:                  sslKey,
				SslCrt:                  sslCrt,
				SslCabundle:             sslCabundle,
			}

			certificate, err = Level27Client.AppSslCertificatesCreateOwn(appID, createOwn)
			if err != nil {
				return err
			}

		case "letsencrypt", "xolphin":
			certificate, err = Level27Client.AppSslCertificatesCreate(appID, create)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("invalid SSL type: %s", appSslCreateSslType)
		}

		if optWait {
			certificate, err = waitForStatus(
				func() (l27.AppSslCertificate, error) {
					return Level27Client.AppSslCertificatesGetSingle(appID, certificate.ID)
				},
				func(s l27.AppSslCertificate) string { return s.Status },
				"ok",
				[]string{"to_create", "creating"},
			)

			if err != nil {
				return fmt.Errorf("waiting on certificate status failed: %s", err.Error())
			}
		}

		outputFormatTemplate(certificate, "templates/entities/appSslCertificate/create.tmpl")
		return nil
	},
}

// APP SSL DELETE
var appSslDeleteForce bool
var appSslDeleteCmd = &cobra.Command{
	Use:     "delete [app] [SSL cert]",
	Short:   "Delete an SSL certificate from an app",
	Example: "lvl app ssl delete forum forum.example.com",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		certID, err := resolveAppSslCertificate(appID, args[1])
		if err != nil {
			return err
		}

		if !appSslDeleteForce {
			app, err := Level27Client.App(appID)
			if err != nil {
				return err
			}

			cert, err := Level27Client.AppSslCertificatesGetSingle(appID, certID)
			if err != nil {
				return err
			}

			if !confirmPrompt(fmt.Sprintf("Delete SSL certificate %s (%d) on app %s (%d)?", cert.Name, certID, app.Name, appID)) {
				return nil
			}
		}

		err = Level27Client.AppSslCertificatesDelete(appID, certID)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/appSslCertificate/delete.tmpl")
		return nil
	},
}

// APP SSL UPDATE
var appSslUpdateCmd = &cobra.Command{
	Use: "update [app] [SSL cert]",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		settings, err := loadMergeSettings(updateSettingsFile, updateSettings)
		if err != nil {
			return err
		}

		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		certID, err := resolveAppSslCertificate(appID, args[1])
		if err != nil {
			return err
		}

		cert, err := Level27Client.AppSslCertificatesGetSingle(appID, certID)
		if err != nil {
			return err
		}

		put := l27.AppSslCertificatePut{
			Name:    cert.Name,
			SslType: cert.SslType,
		}

		data := utils.RoundTripJson(put).(map[string]interface{})
		data = mergeMaps(data, settings)

		err = Level27Client.AppSslCertificatesUpdate(appID, certID, data)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/appSslCertificate/update.tmpl")
		return nil
	},
}

// APP SSL FIX
var appSslFixCmd = &cobra.Command{
	Use:     "fix [app] [SSL cert]",
	Short:   "Fix an invalid SSL certificate",
	Example: "lvl app ssl fix forum forum.example.com",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		certID, err := resolveAppSslCertificate(appID, args[1])
		if err != nil {
			return err
		}

		cert, err := Level27Client.AppSslCertificatesFix(appID, certID)
		if err != nil {
			return err
		}

		outputFormatTemplate(cert, "templates/entities/appSslCertificate/fix.tmpl")
		return nil
	},
}

// APP SSL ACTION
var appSslActionCmd = &cobra.Command{
	Use: "action",
}

// APP SSL ACTION RETRY
var appSslActionRetryCmd = &cobra.Command{
	Use: "retry [app] [SSL cert]",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		certID, err := resolveAppSslCertificate(appID, args[1])
		if err != nil {
			return err
		}

		err = Level27Client.AppSslCertificatesActions(appID, certID, "retry")
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/appSslCertificate/retry.tmpl")
		return nil
	},
}

// APP SSL ACTION VALIDATECHALLENGE
var appSslActionValidateChallengeCmd = &cobra.Command{
	Use: "validateChallenge [app] [SSL cert]",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		certID, err := resolveAppSslCertificate(appID, args[1])
		if err != nil {
			return err
		}

		err = Level27Client.AppSslCertificatesActions(appID, certID, "validateChallenge")
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/appSslCertificate/validateChallenge.tmpl")
		return nil
	},
}

// APP SSL KEY
var appSslKeyCmd = &cobra.Command{
	Use:     "key",
	Short:   "Return a private key for type 'own' sslCertificate.",
	Example: "lvl app ssl key MyAppName",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, err := resolveApp(args[0])
		if err != nil {
			return err
		}

		certID, err := resolveAppSslCertificate(appID, args[1])
		if err != nil {
			return err
		}

		key, err := Level27Client.AppSslCertificatesKey(appID, certID)
		if err != nil {
			return err
		}

		outputFormatTemplate(key, "templates/entities/appSslCertificate/key.tmpl")
		return nil
	},
}
