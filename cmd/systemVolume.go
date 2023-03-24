package cmd

import (
	"fmt"

	"github.com/level27/l27-go"
	"github.com/level27/lvl/utils"
	"github.com/spf13/cobra"
)

func init() {
	// SYSTEM VOLUME
	systemCmd.AddCommand(systemVolumeCmd)

	// SYSTEM VOLUME GET
	systemVolumeCmd.AddCommand(systemVolumeGetCmd)
	addCommonGetFlags(systemVolumeGetCmd)

	// SYSTEM VOLUME CREATE
	systemVolumeCmd.AddCommand(systemVolumeCreateCmd)
	systemVolumeCreateCmd.Flags().StringVar(&systemVolumeCreateName, "name", "", "Name of the new volume")
	systemVolumeCreateCmd.Flags().StringVar(&systemVolumeCreateOrganisation, "organisation", "", "Organisation for the new volume")
	systemVolumeCreateCmd.Flags().StringVar(&systemVolumeCreateDeviceName, "deviceName", "", "Device name for the new volume")
	systemVolumeCreateCmd.Flags().BoolVar(&systemVolumeCreateAutoResize, "autoResize", false, "Enable automatic resizing")
	systemVolumeCreateCmd.Flags().Int32Var(&systemVolumeCreateSpace, "space", 0, "Space of the new volume (in GB)")

	// SYSTEM VOLUME LINK
	systemVolumeCmd.AddCommand(systemVolumeLinkCmd)

	// SYSTEM VOLUME UNLINK
	systemVolumeCmd.AddCommand(systemVolumeUnlinkCmd)

	// SYSTEM VOLUME DELETE
	systemVolumeCmd.AddCommand(systemVolumeDeleteCmd)
	systemVolumeDeleteCmd.Flags().BoolVar(&systemVolumeDeleteForce, "force", false, "Do not ask for confirmation to delete the volume")

	// SYSTEM VOLUME UPDATE
	systemVolumeCmd.AddCommand(systemVolumeUpdateCmd)
	settingsFileFlag(systemVolumeUpdateCmd)
	settingString(systemVolumeUpdateCmd, updateSettings, "name", "New name for the volume")
	settingBool(systemVolumeUpdateCmd, updateSettings, "autoResize", "New autoResize setting")
	settingInt32(systemVolumeUpdateCmd, updateSettings, "space", "New volume space (in GB)")
}

func resolveSystemVolume(systemID l27.IntID, arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
	if err == nil {
		return id, nil
	}

	ip, err := Level27Client.LookupSystemVolumes(systemID, arg)
	if err != nil {
		return 0, err
	}

	if ip == nil {
		return 0, fmt.Errorf("nable to find volume: %s", arg)
	}

	return ip.ID, nil
}

// VOLUMES

// SYSTEM VOLUME
var systemVolumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Commands to manage volumes",
}

// SYSTEM VOLUME GET
var systemVolumeGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get all volumes on a system",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		volumes, err := Level27Client.SystemGetVolumes(systemID, optGetParameters)
		if err != nil {
			return err
		}

		outputFormatTable(
			volumes,
			[]string{"ID", "Name", "Status", "Space", "UID", "AutoResize", "DeviceName"},
			[]string{"ID", "Name", "Status", "Space", "UID", "AutoResize", "DeviceName"})

		return nil
	},
}

// SYSTEM VOLUME CREATE
var systemVolumeCreateName string
var systemVolumeCreateSpace int32
var systemVolumeCreateOrganisation string
var systemVolumeCreateAutoResize bool
var systemVolumeCreateDeviceName string

var systemVolumeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new volume for a system",

	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		organisationID, err := resolveOrganisation(systemVolumeCreateOrganisation)
		if err != nil {
			return err
		}

		create := l27.VolumeCreate{
			Name:         systemVolumeCreateName,
			Space:        systemVolumeCreateSpace,
			Organisation: organisationID,
			System:       systemID,
			AutoResize:   systemVolumeCreateAutoResize,
			DeviceName:   systemVolumeCreateDeviceName,
		}

		volume, err := Level27Client.VolumeCreate(create)
		if err != nil {
			return err
		}

		outputFormatTemplate(volume, "templates/entities/systemVolume/create.tmpl")
		return nil
	},
}

// SYSTEM VOLUME UNLINK
var systemVolumeUnlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "Unlink a volume from a system",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		volumeID, err := resolveSystemVolume(systemID, args[1])
		if err != nil {
			return err
		}

		_, err = Level27Client.VolumeUnlink(volumeID, systemID)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/systemVolume/unlink.tmpl")
		return nil
	},
}

// SYSTEM VOLUME LINK
var systemVolumeLinkCmd = &cobra.Command{
	Use:   "link [system] [volume] [device name]",
	Short: "Link a volume to a system",

	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		// To resolve from name -> ID we need the volume group
		// Easiest way to get that is by getting the volume group ID from the first volume on the system.
		volumes, err := Level27Client.SystemGetVolumes(systemID, l27.CommonGetParams{})
		if err != nil {
			return err
		}

		volumeGroupID := volumes[0].Volumegroup.ID

		volumeID, err := resolveVolumegroupVolume(volumeGroupID, args[1])
		if err != nil {
			return err
		}

		deviceName := args[2]

		_, err = Level27Client.VolumeLink(volumeID, systemID, deviceName)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/systemVolume/link.tmpl")
		return nil
	},
}

// SYSTEM VOLUME DELETE
var systemVolumeDeleteForce bool
var systemVolumeDeleteCmd = &cobra.Command{
	Use:   "delete [system] [volume]",
	Short: "Unlink and delete a volume on a system",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		volumeID, err := resolveSystemVolume(systemID, args[1])
		if err != nil {
			return err
		}

		if !systemVolumeDeleteForce {
			volume, err := Level27Client.VolumeGetSingle(volumeID)
			if err != nil {
				return err
			}

			if !confirmPrompt(fmt.Sprintf("Delete volume %s (%d)?", volume.Name, volume.ID)) {
				return nil
			}
		}

		err = Level27Client.VolumeDelete(volumeID)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/systemVolume/delete.tmpl")
		return nil
	},
}

// SYSTEM VOLUME UPDATE
var systemVolumeUpdateCmd = &cobra.Command{
	Use:   "update [system] [volume]",
	Short: "Update settings on a volume",

	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		settings, err := loadMergeSettings(updateSettingsFile, updateSettings)
		if err != nil {
			return err
		}

		systemID, err := resolveSystem(args[0])
		if err != nil {
			return err
		}

		volumeID, err := resolveSystemVolume(systemID, args[1])
		if err != nil {
			return err
		}

		volume, err := Level27Client.VolumeGetSingle(volumeID)
		if err != nil {
			return err
		}

		volumePut := l27.VolumePut{
			Name:         volume.Name,
			DeviceName:   volume.DeviceName,
			Space:        volume.Space,
			Organisation: volume.Organisation.ID,
			AutoResize:   volume.AutoResize,
			Remarks:      volume.Remarks,
			System:       volume.System.ID,
			Volumegroup:  volume.Volumegroup.ID,
		}

		data := utils.RoundTripJson(volumePut).(map[string]interface{})
		data = mergeMaps(data, settings)

		data["organisation"], err = resolveOrganisation(fmt.Sprint(data["organisation"]))
		if err != nil {
			return err
		}

		err = Level27Client.VolumeUpdate(volumeID, data)
		if err != nil {
			return err
		}

		outputFormatTemplate(nil, "templates/entities/systemVolume/update.tmpl")
		return nil
	},
}
