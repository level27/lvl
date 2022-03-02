package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(regionCommand)

	// region get
	regionCommand.AddCommand(regionGetCommand)

	// region images
	regionCommand.AddCommand(regionImagesCommand)

	// zones
	regionCommand.AddCommand(regionZonesCommand)
}



var regionCommand = &cobra.Command{
	Use: "region",
	Short: "Commands to view available regions for systems",
}

var regionGetCommand = &cobra.Command{
	Use: "get",
	Short: "Get all available regions",

	Run: func(cmd *cobra.Command, args []string) {
		regions := Level27Client.GetRegions()

		outputFormatTable(regions, []string {"ID", "Name", "Country", "Provider"}, []string{"ID", "Name", "Country.Name", "Systemprovider.Name"})
	},
}


var regionImagesCommand = &cobra.Command{
	Use: "images [region]",
	Short: "Get all system images in a region",

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		regionId := regionIdFromArg(args[0])

		regions := Level27Client.GetRegionImages(regionId)

		outputFormatTable(
			regions,
			[]string {"ID", "Name", "OS", "Version"},
			[]string {"ID", "Name", "OperatingsystemVersion.Operatingsystem.Name", "OperatingsystemVersion.Version"})
	},
}

var regionZonesCommand = &cobra.Command{
	Use: "zones",
	Short: "Get all zones in a region",

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		regionId := regionIdFromArg(args[0])

		zones := Level27Client.GetZones(regionId)
		outputFormatTable(zones, []string{"ID", "Name", "Short"}, []string{"ID", "Name", "ShortName"})
	},
}

func regionIdFromArg(arg string) int {
	regionId, err := convertStringToId(arg)
	if err != nil {
		regionMaybe := Level27Client.LookupRegion(arg)
		if regionMaybe == nil {
			log.Fatalln("Unknown region")
		}

		regionId = regionMaybe.ID
	}

	return regionId;
}