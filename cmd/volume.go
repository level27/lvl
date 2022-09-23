package cmd

import (
	"fmt"
	"strconv"

	"github.com/level27/l27-go"
)

func resolveVolumegroupVolume(volumeGroupID int, arg string) (int, error) {
	id, err := strconv.Atoi(arg)
	if err == nil {
		return id, nil
	}

	options, err := Level27Client.LookupVolumegroupVolumes(volumeGroupID, arg)
	if err != nil {
		return 0, err
	}

	res, err := resolveShared(
		options,
		arg,
		"volume",
		func(app l27.Volume) string { return fmt.Sprintf("%s (%d)", app.Name, app.ID) })

	if err != nil {
		return 0, err
	}

	return res.ID, err
}
