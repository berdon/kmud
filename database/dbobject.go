package database

import (
	"labix.org/v2/mgo/bson"
	"sync"
)

type DbObject struct {
	Id bson.ObjectId `bson:"_id"`

	objType   objectType
	mutex     sync.RWMutex
	destroyed bool
}

func (self *DbObject) initDbObject() {
	self.Id = bson.NewObjectId()
}

func (self *DbObject) GetId() bson.ObjectId {
	// Not mutex-protected since thd ID should never change
	return self.Id
}

func (self *DbObject) ReadLock() {
	self.mutex.RLock()
}

func (self *DbObject) ReadUnlock() {
	self.mutex.RUnlock()
}

func (self *DbObject) WriteLock() {
	self.mutex.Lock()
}

func (self *DbObject) WriteUnlock() {
	self.mutex.Unlock()
}

func (self *DbObject) Destroy() {
	self.WriteLock()
	defer self.WriteUnlock()

	self.destroyed = true
}

func (self *DbObject) IsDestroyed() bool {
	self.ReadLock()
	defer self.ReadUnlock()

	return self.destroyed
}

func (self *DbObject) readLocker_str(fn func() string) string {
    self.ReadLock()
    defer self.ReadUnlock()

    return fn()
}

func (self *DbObject) readLocker_id(fn func() bson.ObjectId) bson.ObjectId {
    self.ReadLock()
    defer self.ReadUnlock()

    return fn()
}

func (self *DbObject) readLocker_int(fn func() int) int {
    self.ReadLock()
    defer self.ReadUnlock()

    return fn()
}

func (self *DbObject) writeLocker_str(fn func(string), newVal string, oldVal string) {
    self.WriteLock()
    defer self.WriteUnlock()

    if newVal != oldVal {
        fn(newVal)
    }
}

func (self *DbObject) writeLocker_id(fn func(bson.ObjectId), newVal bson.ObjectId, oldVal bson.ObjectId) {
    self.WriteLock()
    defer self.WriteUnlock()

    if newVal != oldVal {
        fn(newVal)
    }
}

func (self *DbObject) writeLocker_int(fn func(int), newVal int, oldVal int) {
    self.WriteLock()
    defer self.WriteUnlock()

    if newVal != oldVal {
        fn(newVal)
    }
}

// vim: nocindent
