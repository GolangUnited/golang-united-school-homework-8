package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type User struct {
	ID string `json:"id"`
	Email string `json:"email"`
	Age int `json:"age"`	
}

const (
	add = "add"
	list = "list"
	findById = "findById"
	remove = "remove"
)

var (
	ErrOperationMissing = errors.New("-operation flag has to be specified")

)

func parseArgs() Arguments {
	u := User{}
	var(
		OperationFlag string
		ItemFlag string
		FileNameFlag string
	)
	
	flag.StringVar(&OperationFlag, "operation", "list", "takes operations (add, list, findById, remove")
	flag.StringVar(&ItemFlag, "item","", "takes user info")
	flag.StringVar(&FileNameFlag, "fileName","", "tales file name")
	
	
	flag.Parse()
	
	flag.Func("user", "takes user info", func(s string) error {
		dec := json.NewDecoder(strings.NewReader(s))
		err := dec.Decode(&u)
		if err != nil {
			return err
		}
		return nil
	})

	arg := Arguments{
		"id": u.ID,
		"operation": OperationFlag,
		"item": ItemFlag,
		"fileName": FileNameFlag,
	}
	
	return arg
}


type Arguments map[string]string

func Perform(args Arguments, writer io.Writer) error {
	err := CheckErrors(args)
	if err != nil {
		return err
	}
	switch args["operation"]{
	case list:
		err := List(args["fileName"], writer)
		if err != nil {
			return err
		}
	case add:
		err := Add(args)
		if err != nil {
			return err
		}
		
	}
	return nil
}

func Add(arg Arguments) error {
	u := []User{}
	if arg["item"] == ""{
		return fmt.Errorf("-item flag has to be specified")
	}
	f, err := os.OpenFile(arg["fileName"], os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	json.Unmarshal(data, &u)
	for _, r := range u {
		if r.ID == arg["id"] {
			return fmt.Errorf("Item with id %s already exists", r.ID)
		}
	}

	
	return nil
}

func List(fileName string, writer io.Writer) error {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	writer.Write(data)
	
	return nil
}

func CheckErrors(arg Arguments) error {
	if v := arg["operation"]; v==""{
		return ErrOperationMissing
	}else {
		if v != add && v != list && v != findById && v != remove{
			return fmt.Errorf("Operation %v not allowed!", v)
		}
	}
	
	if arg["fileName"] == ""{
		return fmt.Errorf("-fileName flag has to be specified")
	}

	return nil
}

func main() {
	// err := Perform(parseArgs(), os.Stdout)
	// if err != nil {
	// 	panic(err)
	// }

	arg := parseArgs()
	fmt.Println(arg)
}


