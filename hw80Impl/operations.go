package hw80Impl

import (
	"fmt"
	"io"
	"os"
)

func Perform(args Arguments, ioFactory IOFactory, stdOut io.Writer) (err error) {
	if err = args.Validate(); err != nil {
		return err
	}

	var operation string
	if operation, err = args.GetOperation(); err != nil {
		return err
	}

	var strategy performFunc
	var found bool
	if strategy, found = strategies[operation]; !found {
		panic(fmt.Sprintf("operation %s not implemented yet", operation))
	}

	err = strategy(args, ioFactory, stdOut)

	return err
}

type IOFactory struct {
	GetReader func(string) (io.ReadCloser, error)
	GetWriter func(string) (io.WriteCloser, error)
}

func GetOSFactory() IOFactory {
	return IOFactory{
		GetReader: func(fileName string) (io.ReadCloser, error) {
			return os.Open(fileName)
		},
		GetWriter: func(fileName string) (io.WriteCloser, error) {
			return os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_SYNC, 0600)
		},
	}
}

type performFunc func(args Arguments, ioFactory IOFactory, stdOut io.Writer) error

var strategies = map[string]performFunc{
	addOperation:      performAdd,
	listOperation:     performList,
	findByIdOperation: performFindById,
	removeOperation:   performRemove,
}

func performAdd(args Arguments, ioFactory IOFactory, stdOut io.Writer) (err error) {
	var item string
	if item, err = args.GetItem(); err != nil {
		return err
	}

	var fileName string
	if fileName, err = args.GetFileName(); err != nil {
		return err
	}

	itemUser := parseUser(item)

	var users Users
	if users, err = readAllUsers(fileName, ioFactory); err != nil {
		panic(err)
	}

	if users.contains(itemUser.Id) {
		_, _ = fmt.Fprintf(stdOut, itemExistsTemplateError, itemUser.Id)
		return nil
	}

	users = append(users, itemUser)

	if err = users.writeAllUsers(fileName, ioFactory); err != nil {
		panic(err)
	}

	return nil
}

func performList(args Arguments, ioFactory IOFactory, stdOut io.Writer) (err error) {
	var fileName string
	if fileName, err = args.GetFileName(); err != nil {
		return err
	}

	var users Users
	if users, err = readAllUsers(fileName, ioFactory); err != nil {
		panic(err)
	}

	err = users.writeAllJson(stdOut)

	return nil
}

func performFindById(args Arguments, ioFactory IOFactory, stdOut io.Writer) (err error) {
	var id string
	if id, err = args.GetId(); err != nil {
		return err
	}

	var fileName string
	if fileName, err = args.GetFileName(); err != nil {
		return err
	}

	var users Users
	if users, err = readAllUsers(fileName, ioFactory); err != nil {
		return err
	}

	if user, found := users.findById(id); !found {
		return nil
	} else {
		_ = user.writeAllJson(stdOut)
	}

	return nil
}

func performRemove(args Arguments, ioFactory IOFactory, stdOut io.Writer) (err error) {
	var id string
	if id, err = args.GetId(); err != nil {
		return err
	}

	var fileName string
	if fileName, err = args.GetFileName(); err != nil {
		return err
	}

	var users Users
	if users, err = readAllUsers(fileName, ioFactory); err != nil {
		return err
	}

	var usersLen = len(users)
	users = users.filter(func(user User) bool {
		return user.Id != id
	})

	if usersLen == len(users) {
		_, _ = fmt.Fprintf(stdOut, itemNotExistsTemplateError, id)
		return nil
	}

	if err = users.writeAllUsers(fileName, ioFactory); err != nil {
		return err
	}

	return nil
}
