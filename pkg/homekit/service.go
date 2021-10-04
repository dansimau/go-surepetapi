package homekit

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/dansimau/go-surepetapi"
	"github.com/dansimau/go-surepetapi/pkg/config"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
)

type Service struct {
	api *surepetapi.Client
	cfg config.Config

	accessories map[int64]*accessory.Switch
}

func NewService() (*Service, error) {
	cfg, _, err := config.SearchAndLoadConfig("surepet.yaml")
	if err != nil {
		return nil, err
	}

	api, err := surepetapi.NewClient(cfg.API)
	if err != nil {
		return nil, err
	}

	return &Service{
		api: api,
		cfg: *cfg,
	}, nil
}

func (s *Service) detectDevices() {
	res, _, err := s.api.ListDevices()
	if err != nil {
		log.Println(err)
	}

	s.accessories = map[int64]*accessory.Switch{}
	for _, device := range res.Data {
		if device.ProductID == surepetapi.CatFlap {
			log.Println("Adding device", device.ID, device.Name)
			info := accessory.Info{
				ID:           uint64(device.ID),
				Name:         device.Name,
				SerialNumber: device.SerialNumber,
				Manufacturer: "Surepet",
			}
			ac := accessory.NewSwitch(info)

			deviceIDStr := strconv.FormatInt(device.ID, 10)

			ac.Switch.On.OnValueRemoteUpdate(func(on bool) {
				var targetLockState surepetapi.LockState

				if on {
					log.Println("Setting lock for device", device.ID, "to on")
					targetLockState = surepetapi.LockedIn
				} else {
					log.Println("Setting lock for device", device.ID, "to off")
					targetLockState = surepetapi.Unlocked
				}

				done := make(chan struct{})
				go func() {
					_, _, err := s.api.SetLockState(deviceIDStr, targetLockState)
					if err != nil {
						log.Println(err)
					}
					log.Println("Lock state set")
					s.refreshDevices()
					close(done)
				}()

				select {
				case <-done:
				case <-time.After(8 * time.Second):
					log.Println("Timeout; confirming state change anyway")
				}
			})

			s.accessories[device.ID] = ac
		}
	}
}

func (s *Service) refreshDevices() {
	log.Println("Refreshing device status")

	res, _, err := s.api.ListDevices()
	if err != nil {
		log.Println(err)
	}

	for _, device := range res.Data {
		accessory, exists := s.accessories[device.ID]
		if !exists {
			// Add accessory if it's supported
			continue
		}

		// spew.Dump(device.DeviceControl, accessory.Switch.On.GetValue())

		if device.DeviceControl.Locking == surepetapi.Unlocked && accessory.Switch.On.GetValue() {
			accessory.Switch.On.SetValue(false)
			log.Println("Device", device.ID, "has been updated; setting switch to off")
		} else if device.DeviceControl.Locking == surepetapi.LockedIn && !accessory.Switch.On.GetValue() {
			accessory.Switch.On.SetValue(true)
			log.Println("Device", device.ID, "has been updated; setting switch to on")
		}
	}
}

func (s *Service) configureHomekit() (hc.Transport, error) {
	accessoryList := []*accessory.Accessory{}
	for _, accessory := range s.accessories {
		accessoryList = append(accessoryList, accessory.Accessory)
	}

	bridge := accessory.NewBridge(accessory.Info{
		Name: "Surepet Bridge",
	})

	log.Println("Starting homekit transport")

	// configure the ip transport
	config := hc.Config{Pin: "00102003"}
	t, err := hc.NewIPTransport(config, bridge.Accessory, accessoryList...)
	if err != nil {
		return nil, err
	}

	hc.OnTermination(func() {
		log.Println("Shutting down homekit transport")
		<-t.Stop()
	})

	return t, nil
}

func (s *Service) Run() error {
	// Exit on signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		for range sigChan {
			os.Exit(1)
		}
	}()

	s.detectDevices()

	hk, err := s.configureHomekit()
	if err != nil {
		return err
	}

	go hk.Start()

	// Update device status
	s.refreshDevices()
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			s.refreshDevices()
		}
	}()

	select {}
}
