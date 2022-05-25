package cmd

import (
	"fmt"
	"strconv"

	"bitbucket.org/level27/lvl/types"
)

func resolveVolumegroupVolume(volumeGroupID int, arg string) int {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id
	}

	return resolveShared(
		Level27Client.LookupVolumegroupVolumes(volumeGroupID, arg),
		arg,
		"volume",
		func (app types.Volume) string { return fmt.Sprintf("%s (%d)", app.Name, app.ID) }).ID
}