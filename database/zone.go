package database

import (
	"kmud/utils"
)

type Zone struct {
	DbObject `bson:",inline"`

	Name string
}

func NewZone(name string) *Zone {
	var zone Zone
	zone.initDbObject()

	zone.Name = utils.FormatName(name)

	modified(&zone)
	return &zone
}

func (self *Zone) GetType() objectType {
	return ZoneType
}

func (self *Zone) GetName() string {
    return self.readLocker_str(func() string {
        return self.Name
    })
}

func (self *Zone) SetName(name string) {
    self.writeLocker_str(func(name string) {
        self.Name = utils.FormatName(name)
        modified(self)
    }, name, self.Name)
}

type Zones []*Zone

func (self Zones) Contains(z *Zone) bool {
	for _, zone := range self {
		if z == zone {
			return true
		}
	}

	return false
}

// vim: nocindent
