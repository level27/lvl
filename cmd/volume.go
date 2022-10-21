package cmd

import (
	"fmt"

	"github.com/level27/l27-go"
)

func resolveVolumegroupVolume(volumeGroupID l27.IntID, arg string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
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
