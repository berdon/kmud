package database

import (
	"fmt"
	"kmud/utils"
	"labix.org/v2/mgo/bson"
	"sort"
	"strings"
)

type Room struct {
	DbObject `bson:",inline"`

	ZoneId        bson.ObjectId
	Title         string
	Description   string
	Items         []bson.ObjectId
	Location      Coordinate
	ExitNorth     bool
	ExitNorthEast bool
	ExitEast      bool
	ExitSouthEast bool
	ExitSouth     bool
	ExitSouthWest bool
	ExitWest      bool
	ExitNorthWest bool
	ExitUp        bool
	ExitDown      bool
}

type ExitDirection int

const (
	DirectionNorth     ExitDirection = iota
	DirectionNorthEast ExitDirection = iota
	DirectionEast      ExitDirection = iota
	DirectionSouthEast ExitDirection = iota
	DirectionSouth     ExitDirection = iota
	DirectionSouthWest ExitDirection = iota
	DirectionWest      ExitDirection = iota
	DirectionNorthWest ExitDirection = iota
	DirectionUp        ExitDirection = iota
	DirectionDown      ExitDirection = iota
	DirectionNone      ExitDirection = iota
)

type PrintMode int

const (
	ReadMode PrintMode = iota
	EditMode PrintMode = iota
)

func NewRoom(zoneId bson.ObjectId) *Room {
	var room Room
	room.initDbObject()

	room.Title = "The Void"
	room.Description = "You are floating in the blackness of space. Complete darkness surrounds " +
		"you in all directions. There is no escape, there is no hope, just the emptiness. " +
		"You are likely to be eaten by a grue."

	room.ExitNorth = false
	room.ExitNorthEast = false
	room.ExitEast = false
	room.ExitSouthEast = false
	room.ExitSouth = false
	room.ExitSouthWest = false
	room.ExitWest = false
	room.ExitNorthWest = false
	room.ExitUp = false
	room.ExitDown = false

	room.SetLocation(Coordinate{0, 0, 0})
	room.SetZoneId(zoneId)

	modified(&room)
	return &room
}

func (self *Room) GetType() objectType {
	return RoomType
}

func (self *Room) ToString(mode PrintMode, colorMode utils.ColorMode, players []*Character, npcs []*Character, items []*Item) string {
	var str string

	if mode == ReadMode {
		str = fmt.Sprintf("\r\n %v %v %v (%v %v %v)\r\n\r\n %v\r\n\r\n",
			utils.Colorize(colorMode, utils.ColorWhite, ">>>"),
			utils.Colorize(colorMode, utils.ColorBlue, self.GetTitle()),
			utils.Colorize(colorMode, utils.ColorWhite, "<<<"),
			self.GetLocation().X,
			self.GetLocation().Y,
			self.GetLocation().Z,
			utils.Colorize(colorMode, utils.ColorWhite, self.GetDescription()))

		extraNewLine := ""

		if len(players) > 0 {
			str = str + " " + utils.Colorize(colorMode, utils.ColorBlue, "Also here: ")

			var names []string
			for _, char := range players {
				names = append(names, utils.Colorize(colorMode, utils.ColorWhite, char.GetName()))
			}
			str = str + strings.Join(names, utils.Colorize(colorMode, utils.ColorBlue, ", ")) + "\n"

			extraNewLine = "\r\n"
		}

		if len(npcs) > 0 {
			str = str + " " + utils.Colorize(colorMode, utils.ColorBlue, "NPCs: ")

			var names []string
			for _, npc := range npcs {
				names = append(names, utils.Colorize(colorMode, utils.ColorWhite, npc.GetName()))
			}
			str = str + strings.Join(names, utils.Colorize(colorMode, utils.ColorBlue, ", ")) + "\r\n"

			extraNewLine = "\r\n"
		}

		if len(items) > 0 {
			itemMap := make(map[string]int)
			var nameList []string

			for _, item := range items {
				if item == nil {
					continue
				}

				_, found := itemMap[item.GetName()]
				if !found {
					nameList = append(nameList, item.GetName())
				}
				itemMap[item.GetName()]++
			}

			sort.Strings(nameList)

			str = str + " " + utils.Colorize(colorMode, utils.ColorBlue, "Items: ")

			var names []string
			for _, name := range nameList {
				if itemMap[name] > 1 {
					name = fmt.Sprintf("%s x%v", name, itemMap[name])
				}
				names = append(names, utils.Colorize(colorMode, utils.ColorWhite, name))
			}
			str = str + strings.Join(names, utils.Colorize(colorMode, utils.ColorBlue, ", ")) + "\r\n"

			extraNewLine = "\r\n"
		}

		str = str + extraNewLine + " " + utils.Colorize(colorMode, utils.ColorBlue, "Exits: ")

	} else {
		str = fmt.Sprintf(" [1] %v \r\n\r\n [2] %v \r\n\r\n [3] Exits: ", self.GetTitle(), self.GetDescription())
	}

	var exitList []string

	appendIfExists := func(direction ExitDirection) {
		if self.HasExit(direction) {
			exitList = append(exitList, directionToExitString(colorMode, direction))
		}
	}

	appendIfExists(DirectionNorth)
	appendIfExists(DirectionNorthEast)
	appendIfExists(DirectionEast)
	appendIfExists(DirectionSouthEast)
	appendIfExists(DirectionSouth)
	appendIfExists(DirectionSouthWest)
	appendIfExists(DirectionWest)
	appendIfExists(DirectionNorthWest)
	appendIfExists(DirectionUp)
	appendIfExists(DirectionDown)

	if len(exitList) == 0 {
		str = str + utils.Colorize(colorMode, utils.ColorWhite, "None")
	} else {
		str = str + strings.Join(exitList, " ")
	}

	str = str + "\r\n"

	return str
}

func (self *Room) HasExit(dir ExitDirection) bool {
	self.ReadLock()
	defer self.ReadUnlock()

	switch dir {
	case DirectionNorth:
		return self.ExitNorth
	case DirectionNorthEast:
		return self.ExitNorthEast
	case DirectionEast:
		return self.ExitEast
	case DirectionSouthEast:
		return self.ExitSouthEast
	case DirectionSouth:
		return self.ExitSouth
	case DirectionSouthWest:
		return self.ExitSouthWest
	case DirectionWest:
		return self.ExitWest
	case DirectionNorthWest:
		return self.ExitNorthWest
	case DirectionUp:
		return self.ExitUp
	case DirectionDown:
		return self.ExitDown
	}

	panic("Unexpected code path")
}

func (self *Room) SetExitEnabled(dir ExitDirection, enabled bool) {
	self.WriteLock()
	defer self.WriteUnlock()

	switch dir {
	case DirectionNorth:
		self.ExitNorth = enabled
	case DirectionNorthEast:
		self.ExitNorthEast = enabled
	case DirectionEast:
		self.ExitEast = enabled
	case DirectionSouthEast:
		self.ExitSouthEast = enabled
	case DirectionSouth:
		self.ExitSouth = enabled
	case DirectionSouthWest:
		self.ExitSouthWest = enabled
	case DirectionWest:
		self.ExitWest = enabled
	case DirectionNorthWest:
		self.ExitNorthWest = enabled
	case DirectionUp:
		self.ExitUp = enabled
	case DirectionDown:
		self.ExitDown = enabled
	}

	modified(self)
}

func (self *Room) AddItem(item *Item) {
	if !self.HasItem(item) {
		self.WriteLock()
		defer self.WriteUnlock()

		self.Items = append(self.Items, item.GetId())
		modified(self)
	}
}

func (self *Room) RemoveItem(item *Item) {
	if self.HasItem(item) {
		self.WriteLock()
		defer self.WriteUnlock()

		for i, itemId := range self.Items {
			if itemId == item.GetId() {
				// TODO: Potential memory leak. See http://code.google.com/p/go-wiki/wiki/SliceTricks
				self.Items = append(self.Items[:i], self.Items[i+1:]...)
				break
			}
		}

		modified(self)
	}
}

func (self *Room) HasItem(item *Item) bool {
	self.ReadLock()
	defer self.ReadUnlock()

	for _, itemId := range self.Items {
		if itemId == item.GetId() {
			return true
		}
	}

	return false
}

func (self *Room) GetItemIds() []bson.ObjectId {
	self.ReadLock()
	defer self.ReadUnlock()

	return self.Items
}

func (self *Room) SetTitle(title string) {
	self.WriteLock()
	defer self.WriteUnlock()

	if title != self.Title {
		self.Title = title
		modified(self)
	}
}

func (self *Room) GetTitle() string {
	self.ReadLock()
	defer self.ReadUnlock()

	return self.Title
}

func (self *Room) SetDescription(description string) {
	self.WriteLock()
	defer self.WriteUnlock()

	if self.Description != description {
		self.Description = description
		modified(self)
	}
}

func (self *Room) GetDescription() string {
	self.ReadLock()
	defer self.ReadUnlock()

	return self.Description
}

func (self *Room) SetLocation(location Coordinate) {
	self.WriteLock()
	defer self.WriteUnlock()

	if location != self.Location {
		self.Location = location
		modified(self)
	}
}

func (self *Room) GetLocation() Coordinate {
	self.ReadLock()
	defer self.ReadUnlock()

	return self.Location
}

func (self *Room) SetZoneId(zoneId bson.ObjectId) {
	self.WriteLock()
	defer self.WriteUnlock()

	if zoneId != self.ZoneId {
		self.ZoneId = zoneId
		modified(self)
	}
}

func (self *Room) GetZoneId() bson.ObjectId {
	self.ReadLock()
	defer self.ReadUnlock()

	return self.ZoneId
}

func (self *Room) NextLocation(direction ExitDirection) Coordinate {
	loc := self.GetLocation()
	return loc.Next(direction)
}

// vim: nocindent
