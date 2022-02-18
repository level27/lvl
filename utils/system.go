package utils

import (
	"fmt"

	"bitbucket.org/level27/lvl/types"
)

func (c *Client) systemGet(data types.SystemGet) types.SystemGet {

	var result types.SystemGet
	endpoint := fmt.Sprint("/systems")
	c.invokeAPI("GET", endpoint, data, result)
	return result

}
