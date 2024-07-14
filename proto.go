package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/tidwall/resp"
)

const (
	CommandSET = "SET"
	CommandGET = "GET"
)

type Command interface {
}

type SetCommand struct {
	key, val []byte
}

type GetCommand struct {
	key []byte
}

func parseCommand(msg string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(msg))

	for {
		value, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// var cmd Command
		if value.Type() == resp.Array {
			// set command should have 3 elements only
			command := value.Array()
			switch command[0].String() {
			case CommandSET:
				if len(command) != 3 {
					return nil, fmt.Errorf("invalid number of arguments for set : %d", len(command))
				}

				return SetCommand{
					key: command[1].Bytes(),
					val: command[2].Bytes(),
				}, nil

			case CommandGET:
				if len(command) != 2 {
					return nil, fmt.Errorf("invalid number of arguments for set : %d", len(command))
				}

				return GetCommand{
					key: command[1].Bytes(),
				}, nil

			}

		}
	}

	return nil, fmt.Errorf("invalid command received")

}
