package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type user struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// func (u user) MarshalJSON() ([]byte, error) {
// 	j, err := json.Marshal(struct {
// 		Id    string
// 		Email string
// 		Age   int
// 	}{
// 		Id:    u.id,
// 		Email: u.email,
// 		Age:   u.age,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return j, nil
// }

// func (u user) UnmarshalJSON(b []byte) error {
// 	temp := &struct {
// 		Id    string
// 		Email string
// 		Age   int
// 	}{}

// 	if err := json.Unmarshal(b, &temp); err != nil {
// 		return err
// 	}

// 	u.id = temp.Id
// 	u.email = temp.Email
// 	u.age = temp.Age

// 	return nil
// }

type userStorage struct {
	filename string
}

func newUserStorage(filename string) *userStorage {
	storage := userStorage{}
	storage.filename = filename

	return &storage
}

func (s userStorage) getAll() []user {
	byteValue, _ := ioutil.ReadFile(s.filename)
	var users []user
	json.Unmarshal(byteValue, &users)

	return users
}

func (s userStorage) add(u user) error {
	users := s.getAll()
	for _, user := range users {
		if user.Id == u.Id {
			return fmt.Errorf("Item with id %s already exists", u.Id)
		}
	}

	users = append(users, u)

	json, _ := json.Marshal(users)
	_ = ioutil.WriteFile(s.filename, json, 0644)

	return nil
}

func (s userStorage) removeById(id string) error {
	users := s.getAll()
	for i, user := range users {
		if user.Id == id {
			users = append(users[:i], users[i+1:]...)
			json, _ := json.Marshal(users)
			_ = ioutil.WriteFile(s.filename, json, 0644)

			return nil
		}
	}

	return fmt.Errorf("Item with id %s not found", id)
}

func (s userStorage) findById(id string) *user {
	for _, user := range s.getAll() {
		if user.Id == id {
			return &user
		}
	}

	return nil
}
