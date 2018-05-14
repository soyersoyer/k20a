package service

import (
	"github.com/soyersoyer/k20a/db/db"
	"github.com/soyersoyer/k20a/errors"
)

// UserInfoT is struct for clients, stores the user information
type UserInfoT struct {
	Email           string `json:"email"`
	Created         int64  `json:"created"`
	IsAdmin         bool   `json:"is_admin"`
	CollectionCount int    `json:"collection_count"`
}

// GetUsers returns all user
func GetUsers() ([]UserInfoT, error) {
	users, err := db.GetUsers()
	if err != nil {
		return nil, errors.DBError.Wrap(err)
	}
	userInfos := []UserInfoT{}
	for _, u := range users {
		collections, err := db.GetCollectionsOwnedByUser(u.Email)
		if err != nil {
			return nil, errors.DBError.Wrap(err)
		}
		userInfos = append(userInfos, UserInfoT{u.Email, u.Created, u.IsAdmin, len(collections)})
	}
	return userInfos, nil
}

// GetUserInfo fetch an user by the user's email
func GetUserInfo(email string) (*UserInfoT, error) {
	user, err := db.GetUserByEmail(email)
	if err != nil {
		return nil, errors.UserNotExist.T(email).Wrap(err)
	}
	return &UserInfoT{
		user.Email,
		user.Created,
		user.IsAdmin,
		0,
	}, nil
}

// UserUpdateT is the struct for updating a user
type UserUpdateT struct {
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

// UpdateUser updates a user with UserUpdateT struct
func UpdateUser(email string, input *UserUpdateT) error {
	user, err := db.GetUserByEmail(email)
	if err != nil {
		return errors.UserNotExist.T(email).Wrap(err)
	}
	if input.Password != "" {
		hashedPass, err := hashPassword(input.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPass
	}
	user.IsAdmin = input.IsAdmin
	if input.IsAdmin == false {
		admins, err := db.GetAdminUsers()
		if err != nil {
			return err
		}
		if len(admins) == 1 && admins[0].Email == user.Email {
			return errors.AccessDenied
		}
	}
	err = db.UpdateUser(user)
	if err != nil {
		return errors.DBError.Wrap(err)
	}
	return nil
}

// CollectionInfoT is struct for clients, stores the user information
type CollectionInfoT struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	OwnerEmail    string `json:"owner_email"`
	Created       int64  `json:"created"`
	TeammateCount int    `json:"teammate_count"`
}

// GetCollections returns all collection
func GetCollections() ([]CollectionInfoT, error) {
	collections, err := db.GetCollections()
	if err != nil {
		return nil, errors.DBError.Wrap(err)
	}
	collectionInfos := []CollectionInfoT{}
	for _, c := range collections {
		collectionInfos = append(collectionInfos, CollectionInfoT{
			c.ID,
			c.Name,
			c.OwnerEmail,
			c.Created,
			len(c.Teammates)})
	}
	return collectionInfos, nil
}