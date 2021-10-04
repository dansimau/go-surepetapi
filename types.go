package surepetapi

import "time"

type Curfew struct {
	Enabled    bool   `json:"enabled"`
	LockTime   string `json:"lock_time"`
	UnlockTime string `json:"unlock_time"`
}

type Device struct {
	DeviceInfo
	DeviceControl `json:"control"`
	Children      []DeviceInfo `json:"children"`
}

type DeviceControl struct {
	// For cat flaps
	Curfew      []Curfew  `json:"curfew"`
	Locking     LockState `json:"locking"`
	FastPolling bool      `json:"fast_polling"`
}

type DeviceInfo struct {
	ID             int64      `json:"id"`
	ParentDeviceID int64      `json:"parent_device_id"`
	ProductID      EntityType `json:"product_id"`
	HouseholdID    int64      `json:"household_id"`
	Name           string     `json:"name"`
	SerialNumber   string     `json:"serial_number"`
	MACAddress     string     `json:"mac_address"`
	Index          int64      `json:"index"`
	Version        string     `json:"version"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	PairingAt      *time.Time `json:"pairing_at"`
}

type LockState int64

const (
	Unlocked       LockState = 0
	LockedIn       LockState = 1
	LockedOut      LockState = 2
	LockedAll      LockState = 3
	CurfewEnabled  LockState = 4
	CurfewLocked   LockState = -1
	CurfewUnlocked LockState = -2
	CurfewUnknown  LockState = -3
)

type EntityType int64

const (
	Pet        EntityType = 0
	Hub        EntityType = 1
	Repeater   EntityType = 2
	PetFlap    EntityType = 3
	Feeder     EntityType = 4
	Programmer EntityType = 5
	CatFlap    EntityType = 6
	FeederLite EntityType = 7
	Felaqua    EntityType = 8
)

type User struct {
	// 	"id": 2491105121,
	// 	"email_address": "dan@dans.im",
	// 	"first_name": "Daniel",
	// 	"last_name": "Simmons",
	// 	"country_id": 166,
	// 	"language_id": 37,
	// 	"marketing_opt_in": false,
	// 	"terms_accepted": true,
	// 	"weight_units": 0,
	// 	"time_format": 0,
	// 	"version": "MA==",
	// 	"created_at": "2020-09-13T09:01:53+00:00",
	// 	"updated_at": "2020-09-13T09:01:53+00:00",
	// 	"notifications": {
	// 	  "device_status": true,
	// 	  "animal_movement": false,
	// 	  "intruder_movements": false,
	// 	  "new_device_pet": true,
	// 	  "household_management": true,
	// 	  "photos": true,
	// 	  "low_battery": true,
	// 	  "curfew": true,
	// 	  "feeding_activity": true,
	// 	  "drinking_activity": true,
	// 	  "feeding_topup": true,
	// 	  "drinking_topup": true
	// 	}
	//   },
}
