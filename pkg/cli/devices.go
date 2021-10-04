package cli

func init() {
	_, err := parser.AddCommand("devices", "Manage devices", "", &devicesCmd{})
	if err != nil {
		panic(err)
	}
}

type devicesCmd struct {
	List *devicesListCmd `command:"list" alias:"ls" description:"List devices"`
}
