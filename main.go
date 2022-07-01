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

type Arguments map[string]string

type Item struct {
	Id    string
	Email string
	Age   int
}

func Perform(args Arguments, writer io.Writer) error {

	if args["fileName"] == "" {
		return errors.New("-fileName flag has to be specified")
	}

	file, err := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}

	fileToByte, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var jsonItems []Item

	json.Unmarshal(fileToByte, &jsonItems)

	switch args["operation"] {
	case "":
		return errors.New("-operation flag has to be specified")
	case "list":
		listOut, err := json.Marshal(&jsonItems)
		if err == nil {
			listOutPrep := strings.ToLower(string(listOut))
			writer.Write([]byte(listOutPrep))
		}
	case "add":
		if args["item"] != "" {
			var newItem Item
			err := json.Unmarshal([]byte(args["item"]), &newItem)
			if err == nil {
				for _, v := range jsonItems {
					if newItem.Id == v.Id {
						out := fmt.Sprintf("Item with id %s already exists", newItem.Id)
						writer.Write([]byte(out))
					}
				}
				jsonItems = append(jsonItems, newItem)
				out, _ := json.Marshal(jsonItems)

				file.Write([]byte(strings.ToLower(string(out))))
			}
		} else {
			return errors.New("-item flag has to be specified")
		}
	case "findById":
		if args["id"] != "" {
			for i, v := range jsonItems {
				if v.Id == args["id"] {
					out, err := json.Marshal(v)
					if err == nil {
						writer.Write([]byte(strings.ToLower(string(out))))
						return nil
					}
				}
				if i == len(jsonItems) {
					writer.Write([]byte(""))
				}
			}
		} else {
			return errors.New("-id flag has to be specified")
		}
	case "remove":
		if args["id"] != "" {
			for i, v := range jsonItems {
				if v.Id == args["id"] {
					jsonItems[i] = jsonItems[len(jsonItems)-1]
					jsonItems = jsonItems[:len(jsonItems)-1]
					out, err := json.Marshal(jsonItems)
					if err == nil {
						file.Truncate(0)
						file.Seek(0, 0)
						file.Write([]byte(strings.ToLower(string(out))))
						return nil
					}
				}
				if i == len(jsonItems)-1 {
					writer.Write([]byte(fmt.Sprintf("Item with id %s not found", args["id"])))
				}
			}
		} else {
			return errors.New("-id flag has to be specified")
		}
	default:
		return fmt.Errorf("Operation %v not allowed!", args["operation"])
	}

	return nil
}

func parseArgs() Arguments {
	rawId := flag.String("id", "", "Flag to define item's id")
	rawOperation := flag.String("operation", "", "Flag to define operation")
	rawItem := flag.String("item", "", "Flag to define item")
	rawFileName := flag.String("fileName", "", "Flag to define file name")
	flag.Parse()
	toArguments := make(map[string]string)
	toArguments["id"] = *rawId
	toArguments["operation"] = *rawOperation
	toArguments["item"] = *rawItem
	toArguments["fileName"] = *rawFileName
	var args Arguments
	args = toArguments
	return args

}

func main() {

	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		//panic(err)
	}
}
