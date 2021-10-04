package cli

import "strconv"

type devicesListCmd struct{}

func (c *devicesListCmd) Execute(args []string) error {
	api, err := cmd.api()
	if err != nil {
		return err
	}

	devices, _, err := api.ListDevices()
	if err != nil {
		return err
	}

	rows := [][]string{
		{"ID", "Name", "Serial number", "Date created"},
	}
	for _, device := range devices.Data {
		rows = append(rows, []string{
			strconv.FormatInt(device.ID, 10),
			device.Name,
			device.SerialNumber,
			device.CreatedAt.String(),
		})
	}

	printTable(rows)

	return nil
}
