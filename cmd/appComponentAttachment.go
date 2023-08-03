package cmd

import (
	b64 "encoding/base64"
	"fmt"
	"io"

	"github.com/level27/l27-go"
	"github.com/spf13/cobra"
)

func init() {
	appComponentCmd.AddCommand(appComponentAttachmentCmd)

	appComponentAttachmentCmd.AddCommand(appComponentAttachmentUploadCmd)
	appComponentAttachmentUploadCmd.Flags().StringVar(
		&appComponentAttachmentUploadOrganisation, "organisation", "", "Organisation to upload this attachment for")
	appComponentAttachmentUploadCmd.Flags().StringVarP(
		&appComponentAttachmentUploadName, "name", "n", "", "Name of the created attachment")
	appComponentAttachmentUploadCmd.Flags().StringVarP(
		&appComponentAttachmentUploadType, "type", "t", "", "Type of the attachment")
}

var appComponentAttachmentCmd = &cobra.Command{
	Use:     "attachment",
	Aliases: []string{"attachments"},
	Short:   "Commands for managing app component attachments",
	Long: `Some components (such as solr) allow you to associate a bundle of files uploaded in advance. This is done through 'attachments'.
These commands allow you to upload & manage these attachments. To actually use them with a component, specify --attachment to the relevant lvl command, such as lvl app component create.
`,
	Example: `To upload & use a solr attachment:
lvl app component attachment upload ./config.zip -t solr -n solr_config
lvl app component create my_app --type solr --name solr --system my.cool.system --attachment <id>`,
}

var appComponentAttachmentUploadOrganisation string
var appComponentAttachmentUploadName string
var appComponentAttachmentUploadType string

var appComponentAttachmentUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a new attachment",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		org, err := resolveOrgOrUserOrg(appComponentAttachmentUploadOrganisation)
		if err != nil {
			return fmt.Errorf("couldn't resolve organisation: %v", err)
		}

		data, err := func() ([]byte, error) {
			file, err := openArgFile(args[0])
			if err != nil {
				return nil, err
			}

			defer file.Close()

			return io.ReadAll(file)
		}()

		if err != nil {
			return fmt.Errorf("couldn't read file: %v", err)
		}

		encoded := b64.StdEncoding.EncodeToString(data)

		upload := l27.AttachmentUpload{
			Name:         appComponentAttachmentUploadName,
			Type:         &appComponentAttachmentUploadType,
			EntityClass:  "Level27\\AppBundle\\Entity\\Appcomponent",
			Organisation: org,
			File:         encoded,
		}

		attachment, err := Level27Client.AttachmentUpload(upload)
		if err != nil {
			return fmt.Errorf("API error: %v", err)
		}

		outputFormatTemplate(attachment, "templates/entities/appComponentAttachment/upload.tmpl")
		return nil
	},
}
