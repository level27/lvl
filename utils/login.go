package utils

import (
	"fmt"

	"bitbucket.org/level27/lvl/types"
)

func (c *Client) Login(username string, password string) types.Login {

	var login types.Login
	endpoint := "login"
	data := fmt.Sprintf("{\"username\":\"%s\", \"password\":\"%s\"}", username, password)
	c.invokeAPI("POST", endpoint, data, &login)

	return login
}
