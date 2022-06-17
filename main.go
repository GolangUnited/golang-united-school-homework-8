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
	
	var(
		OperationFlag string
		ItemFlag string
		FileNameFlag string
		IdFlag string
	)
	
	flag.StringVar(&OperationFlag, "operation", "list", "takes operations (add, list, findById, remove)")
	flag.StringVar(&ItemFlag, "item","", "takes user info")
	flag.StringVar(&FileNameFlag, "fileName","users.json", "tales file name")
	flag.StringVar(&IdFlag, "id", "", "takes an id")
	flag.Parse()
	
	
	arg := Arguments{
		"id": IdFlag,
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
		err := Add(args, writer)

		if err != nil {
			return err
		}
	case remove:
		err := Remove(args, writer)
		if err != nil{
			return err
		}
	case findById:
		err := FindById(args, writer)
		if err != nil {
			return err
		}
		
	}
	return nil
}

func FindById(arg Arguments, w io.Writer) error {
	content := []User{}
	if arg["id"] == "" {
		return fmt.Errorf("-id flag has to be specified")
	}
	f, err := os.Open(arg["fileName"])
	if err != nil {
		return err
	}	
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	json.Unmarshal(data, &content)
	for _, d := range content {
		if d.ID == arg["id"] {
			str := fmt.Sprintf("{\"id\":\"%s\",\"email\":\"%s\",\"age\":%v}", d.ID, d.Email, d.Age)
			w.Write([]byte(str))
			return nil
			
		}
	}
	w.Write([]byte(""))
	return nil
}

func Remove(arg Arguments, w io.Writer) error {
	dataUser := []User{}
	if arg["id"] == "" {
		return fmt.Errorf("-id flag has to be specified")
	}

	f, err := os.Open(arg["fileName"])
	if err != nil {
		return err
	}	

	exixstingData, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	json.Unmarshal(exixstingData, &dataUser)
	
	inx := 0
	exist := false
	for i, u := range dataUser {
		if u.ID == arg["id"]{
			exist = true
			inx = i
		}
	}

	if !exist {
		message := fmt.Sprintf("Item with id %s not found", arg["id"])
		w.Write([]byte(message))
		return nil
	}

	os.Remove(arg["fileName"])
	newF, _ := os.Create(arg["fileName"])
	dataUser = append(dataUser[:inx],dataUser[inx+1:]... )
	dataForFile, _ := json.Marshal(dataUser)
	newF.Write(dataForFile)

	return nil
}

func Add(arg Arguments, w io.Writer) error {
	fu := []User{}
	newUser := User{}
	newData := []byte(arg["item"])
	json.Unmarshal(newData, &newUser)

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
		if newUser.ID == r.ID{
			message := fmt.Sprintf("Item with id %s already exists", r.ID)
			w.Write([]byte(message))
			return nil 
		}
	}
	os.Remove(arg["fileName"])
	newF, _ := os.Create(arg["fileName"])
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
}


