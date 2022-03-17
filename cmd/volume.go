package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func resolveVolumegroupVolume(volumeGroupID int, arg string) int {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id
	}

	volume := Level27Client.LookupVolumegroupVolumes(volumeGroupID, arg)
	if volume == nil {
		cobra.CheckErr(fmt.Sprintf("Unable to find volume: %s", arg))
		return 0
	}

	return volume.ID
}