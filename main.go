package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	
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
type Arguments map[string]string


func parseArgs() Arguments {
	u := User{}
	var(
		OperationFlag string
		ItemFlag string
		FileNameFlag string
	)
	
	flag.StringVar(&OperationFlag, "operation", "list", "takes operations (add, list, findById, remove)")
	flag.StringVar(&ItemFlag, "item","", "takes user info")
	flag.StringVar(&FileNameFlag, "fileName","users.json", "tales file name")
	
	flag.Parse()
	
	data := []byte(ItemFlag)
	err := json.Unmarshal(data, &u)
	if err != nil {
		fmt.Println(err)
	}
	arg := Arguments{
		"id": u.ID,
		"operation": OperationFlag,
		"item": ItemFlag,
		"fileName": FileNameFlag,
	}
	
	return arg
}


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
		message := fmt.Sprintf("Item with id %s already exists", args["id"])
		if err == errors.New(message) {
			writer.Write([]byte(message))
		}
		if err != nil {
			return err
		}
	case remove:
		
	}
	return nil
}

func Add(arg Arguments) error {
	fu := []User{}

	if arg["item"] == ""{
		return fmt.Errorf("-item flag has to be specified")
	}

	f, err := os.OpenFile(arg["fileName"], os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	json.Unmarshal(data, &fu)
	

	for _, r := range fu {
		fmt.Println("arg ID: ", arg["id"], "r.ID: ", r.ID)
		if arg["id"] == r.ID{
			return fmt.Errorf("Item with id %s already exists", r.ID)
		}
	}
	os.Remove(arg["fileName"])
	newF, _ := os.Create(arg["fileName"])
	newUser := User{}
	newData := []byte(arg["item"])
	json.Unmarshal(newData, &newUser)
	fu = append(fu, newUser)

	dataForFile, _ := json.Marshal(fu)
	newF.Write(dataForFile)
	
	
	return nil
}

func List(fileName string, writer io.Writer) error {
	f, err := os.Open(fileName)
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

func Remove()

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
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}

	// arg := parseArgs()
	// fmt.Println(arg)
	// err := Add(arg)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// List(arg["fileName"], os.Stdout)

}


