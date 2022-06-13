package hw80Impl

import (
	"encoding/json"
	"io"
	"strings"
)

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type Users []User

func parseUsers(data []byte) (result Users, err error) {
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func parseUser(json string) User {
	if users, err := parseUsers([]byte("[" + strings.Trim(json, "[]") + "]")); err != nil {
		panic(json)
	} else if len(users) < 1 {
		panic(json)
	} else {
		return users[0]
	}
}

func (receiver Users) toBytes() (result []byte, err error) {
	if result, err = json.Marshal(&receiver); err != nil {
		return nil, err
	}

	return result, nil
}

func (receiver Users) findById(id string) (User, bool) {
	for _, user := range receiver {
		if user.Id == id {
			return user, true
		}
	}

	return User{}, false
}

func (receiver Users) contains(id string) bool {
	_, found := receiver.findById(id)

	return found
}

func (receiver Users) writeAllUsers(fileName string, ioFactory IOFactory) (err error) {
	var writer io.WriteCloser
	defer func() {
		_ = writer.Close()
	}()

	var data []byte
	if data, err = receiver.toBytes(); err != nil {
		panic(err)
	}

	if writer, err = ioFactory.GetWriter(fileName); err != nil {
		return err
	}

	if _, err = writer.Write(data); err != nil {
		panic(err)
	}

	return nil
}

func (receiver Users) writeAllJson(stdOut io.Writer) (err error) {
	if len(receiver) < 1 {
		return nil
	}

	var data []byte
	if data, err = json.Marshal(receiver); err != nil {
		return err
	}

	if _, err = stdOut.Write(data); err != nil {
		return err
	}

	return nil
}

func (receiver Users) filter(selectFilter func(user User) bool) (result Users) {
	result = make(Users, 0, len(receiver))
	for _, user := range receiver {
		if !selectFilter(user) {
			continue
		}

		result = append(result, user)
	}

	return result
}

func (receiver User) writeAllJson(stdOut io.Writer) (err error) {
	var data []byte
	if data, err = json.Marshal(receiver); err != nil {
		return err
	}

	if _, err = stdOut.Write(data); err != nil {
		return err
	}

	return nil
}

func readAllUsers(fileName string, ioFactory IOFactory) (users Users, err error) {
	var reader io.ReadCloser
	defer func() {
		_ = reader.Close()
	}()

	if reader, err = ioFactory.GetReader(fileName); err != nil {
		return Users{}, nil
	}

	var data []byte
	if data, err = io.ReadAll(reader); err != nil {
		return nil, err
	}

	if users, err = parseUsers(data); err != nil {
		return nil, err
	}

	return users, nil
}
