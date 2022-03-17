package utils

import (
	"fmt"

	"bitbucket.org/level27/lvl/types"
)

// GET /volume/{volumeID}
func (c *Client) VolumeGetSingle(volumeID int) types.Volume {
	var response struct {
		Volume types.Volume `json:"volume"`
	}

	endpoint := fmt.Sprintf("volumes/%d", volumeID)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "VolumeGetSingle")

	return response.Volume
}

// GET /volume
func (c *Client) VolumeGetList(get types.CommonGetParams) []types.Volume {
	var response struct {
		Volumes []types.Volume `json:"volumes"`
	}

	endpoint := "volumes"
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "VolumeGetList")

	return response.Volumes
}

// POST /volume
func (c *Client) VolumeCreate(create types.VolumeCreate) types.Volume {
	var response struct {
		Volume types.Volume `json:"volume"`
	}

	endpoint := "volumes"
	err := c.invokeAPI("POST", endpoint, create, &response)
	AssertApiError(err, "VolumeCreate")

	return response.Volume
}

// DELETE /volume/{volumeID}
func (c *Client) VolumeDelete(volumeID int) {
	endpoint := fmt.Sprintf("volumes/%d", volumeID)
	err := c.invokeAPI("DELETE", endpoint, nil, nil)
	AssertApiError(err, "VolumeCreate")
}

// PUT /volume/{volumeID}
func (c *Client) VolumeUpdate(volumeID int, data map[string]interface{}) {
	endpoint := fmt.Sprintf("volumes/%d", volumeID)
	err := c.invokeAPI("PUT", endpoint, data, nil)
	AssertApiError(err, "VolumeUpdate")
}

// POST /volume/{volumeID}/actions (link)
func (c *Client) VolumeLink(volumeID int, systemID int, deviceName string) types.Volume {
	var response struct {
		Volume types.Volume `json:"volume"`
	}

	var request struct {
		Type string `json:"type"`
		System int `json:"system"`
		DeviceName string `json:"deviceName"`
	}

	request.Type = "link"
	request.System = systemID
	request.DeviceName = deviceName

	endpoint := fmt.Sprintf("volumes/%d/actions", volumeID)
	err := c.invokeAPI("POST", endpoint, request, &response)
	AssertApiError(err, "VolumeLink")

	return response.Volume
}

// POST /volume/{volumeID}/actions (unlink)
func (c *Client) VolumeUnlink(volumeID int, systemID int) types.Volume {
	var response struct {
		Volume types.Volume `json:"volume"`
	}

	var request struct {
		Type string `json:"type"`
		System int `json:"system"`
	}

	request.Type = "unlink"
	request.System = systemID

	endpoint := fmt.Sprintf("volumes/%d/actions", volumeID)
	err := c.invokeAPI("POST", endpoint, request, &response)
	AssertApiError(err, "VolumeUnlink")

	return response.Volume
}

// GET /volumegroups/{volumegroupID}/volumes
func (c *Client) VolumegroupVolumeGetList(volumegroupID int, get types.CommonGetParams) []types.Volume {
	var response struct {
		Volumes []types.Volume `json:"volumes"`
	}

	endpoint := fmt.Sprintf("volumegroups/%d/volumes?%s", volumegroupID, formatCommonGetParams(get))
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "VolumegroupVolumeGetList")

	return response.Volumes
}

func (c *Client) LookupVolumegroupVolumes(volumeGroupID int, name string) *types.Volume {
	volumes := c.VolumegroupVolumeGetList(volumeGroupID, types.CommonGetParams{Filter: name})
	for _, volume := range volumes {
		if volume.Name == name {
			return &volume
		}
	}

	return nil
}