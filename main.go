package main

import (
	"flag"
	"io"
	"os"
)

type Arguments map[string]string

func Perform(args Arguments, writer io.Writer) error {
	operation := configureOperations()
	return operation.execute(args, writer)
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

func configureOperations() operation {
	validation := &argumentsValidationOperation{}

	list := &listOperation{}
	validation.setNext(list)

	add := &addOperation{}
	list.setNext(add)

	findById := &findByIdOperation{}
	add.setNext(findById)

	remove := &removeOperation{}
	findById.setNext(remove)

	incorrectOp := &incorrectOperationCase{}
	remove.setNext(incorrectOp)

	return validation
}

func parseArgs() Arguments {
	args := Arguments{}
	for _, argKey := range []string{operationArgKey, filenameArgKey, itemArgKey, idArgKey} {
		flag.String(args[argKey], argKey, "")
	}
	flag.Parse()

	return args
}
