package cli

import (
	"fmt"

	"github.com/dansimau/go-surepetapi"
)

func init() {
	_, err := parser.AddCommand("lock", "Lock a catflap", "", &lockCmd{})
	if err != nil {
		panic(err)
	}
}

type lockCmd struct {
	Arguments struct {
		DeviceID string
	} `positional-args:"true" required:"true"`
}

func (c *lockCmd) Execute(args []string) error {
	api, err := cmd.api()
	if err != nil {
		return err
	}

	res, _, err := api.SetLockState(c.Arguments.DeviceID, surepetapi.LockedIn)
	if err != nil {
		return err
	}

	fmt.Println(res.Locking)

	return nil
}
