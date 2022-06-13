package hw80Impl

import (
	"flag"
	"fmt"
	"os"
)

const (
	operationFlag = "operation"
	fileNameFlag  = "fileName"
	itemFlag      = "item"
	idFlag        = "id"

	addOperation      = "add"
	listOperation     = "list"
	findByIdOperation = "findById"
	removeOperation   = "remove"

	operationFlagMissingError   = "-operation flag has to be specified"
	wrongOperationTemplateError = "Operation %s not allowed!"
	idFlagMissingError          = "-id flag has to be specified"
	fileNameMissingError        = "-fileName flag has to be specified"
	itemFlagMissingError        = "-item flag has to be specified"

	itemExistsTemplateError    = "Item with id %s already exists"
	itemNotExistsTemplateError = "Item with id %s not found"
)

var (
	supportedFlags = [...]struct {
		name         string
		defaultValue string
		usage        string
	}{{name: operationFlag}, {name: fileNameFlag}, {name: itemFlag}, {name: idFlag}}
	supportedOperations = strArr{addOperation, listOperation, findByIdOperation, removeOperation}
)

type Arguments map[string]string
type strArr []string

func (r *Arguments) GetOperation() (string, error) {
	return r.getArgValue(operationFlag, operationFlagMissingError)
}

func (r *Arguments) GetId() (string, error) {
	return r.getArgValue(idFlag, idFlagMissingError)
}

func (r *Arguments) GetFileName() (string, error) {
	return r.getArgValue(fileNameFlag, fileNameMissingError)
}

func (r *Arguments) GetItem() (string, error) {
	return r.getArgValue(itemFlag, itemFlagMissingError)
}

func (r *Arguments) getArgValue(key string, errorMessage string) (string, error) {
	if result, ok := (*r)[key]; !ok || result == "" {
		return "", fmt.Errorf(errorMessage)
	} else {
		return result, nil
	}
}

func parseFlagSet(flagSet *flag.FlagSet, args []string) Arguments {
	for _, item := range supportedFlags {
		flagSet.String(item.name, item.defaultValue, item.usage)
	}

	if err := flagSet.Parse(args); err != nil {
		panic(err)
	}

	result := Arguments{}
	flagSet.Visit(func(f *flag.Flag) {
		result[f.Name] = f.Value.String()
	})

	return result
}

func (arr strArr) contains(item string) bool {
	for _, s := range arr {
		if s == item {
			return true
		}
	}

	return false
}

func (r *Arguments) Validate() error {
	var operation string
	var err error
	if operation, err = r.GetOperation(); err != nil {
		return err
	}

	if !supportedOperations.contains(operation) {
		//goland:noinspection GoErrorStringFormat
		return fmt.Errorf(wrongOperationTemplateError, operation)
	}

	if _, err = r.GetFileName(); err != nil {
		return err
	}

	return nil
}

func NewArgumentsFromCommandLine() Arguments {
	return parseFlagSet(flag.CommandLine, os.Args[1:])
}
