package utils

import (
	"fmt"

	"bitbucket.org/level27/lvl/types"
)

func (c *Client) GetRegions() []types.Region {
	var response struct {
		Regions []types.Region `json:"regions"`
	}

	err := c.invokeAPI("GET", "regions", nil, &response)
	AssertApiError(err, "GetRegions")

	return response.Regions
}

// Try to get a region by name
func (c *Client) LookupRegion(name string) *types.Region {
	regions := c.GetRegions()
	for _, region := range regions {
		if region.Name == name {
			return &region
		}
	}

	return nil
}

// Try to get a zone by name.
// Very slow.
func (c *Client) LookupZone(name string) *types.Zone {
	regions := c.GetRegions()
	for _, region := range regions {
		for _, zone := range c.GetZones(region.ID) {
			if zone.Name == name {
				return &zone
			}
		}
	}

	return nil
}

func (c *Client) GetZones(region int) []types.Zone {
	var response struct {
		Zones []types.Zone `json:"zones"`
	}

	endpoint := fmt.Sprintf("regions/%d/zones", region)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "GetZones")

	return response.Zones
}

func (c *Client) GetRegionImages(region int) []types.Image {
	var response struct {
		Images []types.Image `json:"systemimages"`
	}

	endpoint := fmt.Sprintf("regions/%d/images", region)
	err := c.invokeAPI("GET", endpoint, nil, &response)
	AssertApiError(err, "GetRegions")

	return response.Images
}