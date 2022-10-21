package cmd

import (
	"fmt"
	"strings"

	"github.com/level27/l27-go"
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

func resolveZoneRegion(zoneName string) (l27.IntID, l27.IntID, error) {
	zone, region, err := Level27Client.LookupZoneAndRegion(zoneName)
	if err != nil {
		return 0, 0, err
	}

	if zone == nil || region == nil {
		return 0, 0, fmt.Errorf("unable to find zone: %s", zoneName)
	}

	return zone.ID, region.ID, nil
}

func resolveRegionImage(region l27.IntID, imageName string) (l27.IntID, error) {
	id, err := l27.ParseID(imageName)
	if err == nil {
		return id, nil
	}

	splittedImageData := strings.SplitN(imageName, " ", 2)
	osName := splittedImageData[0]
	osVersion := splittedImageData[1]

	images, err := Level27Client.GetRegionImages(region)
	if err != nil {
		return 0, err
	}

	for _, image := range images {
		if image.OperatingsystemVersion.Version == osVersion && image.OperatingsystemVersion.Operatingsystem.Name == osName {
			return image.ID, nil
		}
	}

	return 0, fmt.Errorf("unable to find image with name %s in zone", imageName)
}

var regionCommand = &cobra.Command{
	Use:   "region",
	Short: "Commands to view available regions for systems",
}

var regionGetCommand = &cobra.Command{
	Use:   "get",
	Short: "Get all available regions",

	RunE: func(cmd *cobra.Command, args []string) error {
		regions, err := Level27Client.GetRegions()
		if err != nil {
			return err
		}

		outputFormatTable(regions, []string{"ID", "Name", "Country", "Provider"}, []string{"ID", "Name", "Country.Name", "Systemprovider.Name"})
		return nil
	},
}

var regionImagesCommand = &cobra.Command{
	Use:   "images [region]",
	Short: "Get all system images in a region",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		regionID, err := regionIDFromArg(args[0])
		if err != nil {
			return err
		}

		regions, err := Level27Client.GetRegionImages(regionID)
		if err != nil {
			return err
		}

		outputFormatTable(
			regions,
			[]string{"ID", "Name", "OS", "Version"},
			[]string{"ID", "Name", "OperatingsystemVersion.Operatingsystem.Name", "OperatingsystemVersion.Version"})

		return nil
	},
}

var regionZonesCommand = &cobra.Command{
	Use:   "zones",
	Short: "Get all zones in a region",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		regionID, err := regionIDFromArg(args[0])
		if err != nil {
			return err
		}

		zones, err := Level27Client.GetZones(regionID)
		if err != nil {
			return err
		}

		outputFormatTable(zones, []string{"ID", "Name", "Short"}, []string{"ID", "Name", "ShortName"})
		return nil
	},
}

func regionIDFromArg(arg string) (l27.IntID, error) {
	regionID, err := convertStringToID(arg)
	if err != nil {
		regionMaybe, err := Level27Client.LookupRegion(arg)
		if err != nil {
			return 0, err
		}

		if regionMaybe == nil {
			return 0, fmt.Errorf("unknown region: '%s'", arg)
		}

		regionID = regionMaybe.ID
	}

	return regionID, nil
}
