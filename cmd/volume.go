package cmd

import (
	"fmt"
	"strconv"

	"github.com/level27/l27-go"
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
		func(app l27.Volume) string { return fmt.Sprintf("%s (%d)", app.Name, app.ID) }).ID
}
