package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var (
	operationArgKey string = "operation"
	filenameArgKey  string = "fileName"
	itemArgKey      string = "item"
	idArgKey        string = "id"
)

type operation interface {
	execute(Arguments, io.Writer) error
	setNext(operation)
}

//add list findById remove
type argumentsValidationOperation struct {
	next operation
}

type addOperation struct {
	next operation
}

type listOperation struct {
	next operation
}

type findByIdOperation struct {
	next operation
}

type removeOperation struct {
	next operation
}

type incorrectOperationCase struct {
	next operation
}

func (op *argumentsValidationOperation) execute(args Arguments, writer io.Writer) error {
	if value, ok := args[operationArgKey]; !ok || len(value) < 1 {
		return errors.New("-operation flag has to be specified")
	}

	if value, ok := args[filenameArgKey]; !ok || len(value) < 1 {
		return errors.New("-fileName flag has to be specified")
	}

	return op.next.execute(args, writer)
}

func (op *addOperation) execute(args Arguments, writer io.Writer) error {
	if args[operationArgKey] != "add" {
		return op.next.execute(args, writer)
	}

	if value, ok := args[itemArgKey]; !ok || len(value) < 1 {
		return errors.New("-item flag has to be specified")
	}

	var user user
	json.Unmarshal([]byte(args[itemArgKey]), &user)

	storage := newUserStorage(args[filenameArgKey])
	err := storage.add(user)

	if err != nil {
		writer.Write([]byte(err.Error()))
	}

	return nil
}

func (op *listOperation) execute(args Arguments, writer io.Writer) error {
	if args[operationArgKey] != "list" {
		return op.next.execute(args, writer)
	}

	storage := newUserStorage(args[filenameArgKey])
	list := storage.getAll()
	json, _ := json.Marshal(list)
	writer.Write(json)

	return nil
}

func (op *findByIdOperation) execute(args Arguments, writer io.Writer) error {
	if args[operationArgKey] != "findById" {
		return op.next.execute(args, writer)
	}

	if value, ok := args[idArgKey]; !ok || len(value) < 1 {
		return errors.New("-id flag has to be specified")
	}

	storage := newUserStorage(args[filenameArgKey])
	user := storage.findById(args[idArgKey])
	if user != nil {
		json, _ := json.Marshal(user)
		writer.Write(json)
	} else {
		writer.Write([]byte{})
	}

	return nil
}

func (op *removeOperation) execute(args Arguments, writer io.Writer) error {
	if args[operationArgKey] != "remove" {
		return op.next.execute(args, writer)
	}

	if value, ok := args[idArgKey]; !ok || len(value) < 1 {
		return errors.New("-id flag has to be specified")
	}

	storage := newUserStorage(args[filenameArgKey])
	err := storage.removeById(args[idArgKey])

	if err != nil {
		writer.Write([]byte(err.Error()))
	}

	return nil
}

func (op *incorrectOperationCase) execute(args Arguments, writer io.Writer) error {
	return fmt.Errorf("Operation %s not allowed!", args[operationArgKey])
}

func (op *argumentsValidationOperation) setNext(next operation) {
	op.next = next
}

func (op *addOperation) setNext(next operation) {
	op.next = next
}

func (op *listOperation) setNext(next operation) {
	op.next = next
}

func (op *findByIdOperation) setNext(next operation) {
	op.next = next
}

func (op *removeOperation) setNext(next operation) {
	op.next = next
}

func (op *incorrectOperationCase) setNext(next operation) {
	op.next = next
}
