package homekit

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type Doorlock struct {
	*accessory.Accessory
	LockManagement *service.LockManagement
}

func NewDoorlock(info accessory.Info) *Doorlock {
	acc := Doorlock{}
	acc.Accessory = accessory.New(info, accessory.TypeLightbulb)
	acc.LockManagement = service.NewLockManagement()

	acc.AddService(acc.LockManagement.Service)

	return &acc
}
