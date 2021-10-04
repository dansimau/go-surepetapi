package cli

import (
	"fmt"

	"github.com/dansimau/go-surepetapi"
)

func init() {
	_, err := parser.AddCommand("unlock", "Unlock a catflap", "", &unlockCmd{})
	if err != nil {
		panic(err)
	}
}

type unlockCmd struct {
	Arguments struct {
		DeviceID string
	} `positional-args:"true" required:"true"`
}

func (c *unlockCmd) Execute(args []string) error {
	api, err := cmd.api()
	if err != nil {
		return err
	}

	res, _, err := api.SetLockState(c.Arguments.DeviceID, surepetapi.Unlocked)
	if err != nil {
		return err
	}

	fmt.Println(res.Locking)

	return nil
}
