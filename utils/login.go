package utils

import (
	"bitbucket.org/level27/lvl/types"
)

func (c *Client) Login(username string, password string) (types.Login, error) {
	var login types.Login

	err := c.invokeAPI("POST", "login", &types.LoginRequest{Username: username, Password: password}, &login)

	return login, err
}
