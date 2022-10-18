package unit21

import (
	"encoding/json"
	"log"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
)

type Entity interface {
	Create(user model.User) (unit21Id string, err error)
	Update(id string, updates any) (err error)
}

type entity struct {
	userRepo repository.User
	// deviceRepo     repository.Device
	contactRepo repository.UserContact
	// instrumentRepo repository.Instrument
}

// With Device and Instrument
// func newEntity(user repository.User, device repository.Device, contact repository.UserContact, instrument repository.Instrument) Entity {
// 	return &entity{userRepo: user, deviceRepo: device, contactRepo: contact, instrumentRepo: instrument}
// }

func newEntity(user repository.User, contact repository.UserContact) Entity {
	return &entity{userRepo: user, contactRepo: contact}
}

// //
// u21entity := unit21.NewEntity()
// u21insturment := unit21.NewIntrument()

// u21entity.Create(user)
// //

func (e entity) Create(user model.User) (unit21Id string, err error) {

	// ultimately may want a join here.
	// devices, err  := e.deviceRepo.ListUserID(user.ID, 100, 0)
	// instruments, err  := e.instrumentRepo.ListUserID(user.ID, 100, 0)
	_, err = e.contactRepo.ListUserID(user.ID, 100, 0)

	body, err := create("entities", MapUserToEntity(user))

	var entity *entityResponse
	err = json.Unmarshal(body, &entity)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return
	}

	log.Printf("Unit21Id: %s", entity.Unit21Id)
	return entity.Unit21Id, nil
}

func (e entity) Update(id string, updates any) (err error) {

	_, err = update("entities", id, updates)

	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return
	}

	return
}
