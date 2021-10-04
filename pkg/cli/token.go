package cli

import (
	"fmt"

	"github.com/dansimau/go-surepetapi"
)

func init() {
	_, err := parser.AddCommand("token", "Generate token", "", &tokenCmd{})
	if err != nil {
		panic(err)
	}
}

type tokenCmd struct {
	EmailAddress string `long:"email-address" description:"Email address associated with the Surepet account"`
	Password     string `long:"password" description:"Password for the Surepet account"`
	DeviceID     string `long:"device-id" description:"A unique ID for this device"`
}

func (c *tokenCmd) Execute(args []string) error {
	api, err := cmd.api()
	if err != nil {
		return err
	}

	res, _, err := api.AuthLogin(surepetapi.AuthLoginRequest{
		EmailAddress: c.EmailAddress,
		Password:     c.Password,
		DeviceID:     c.DeviceID,
	})
	if err != nil {
		return err
	}

	fmt.Println(res.Data.Token)

	return nil
}
