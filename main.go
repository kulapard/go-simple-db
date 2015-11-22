package main


import (
	"fmt"
	"strconv"
	"strings"
	"bufio"
	"os"
	"errors"
)


var dbStorage = map[string]string{}
var transactions = NestedTransactions{}

// Errors
var ErrNoTransaction = errors.New("NO TRANSACTION")
var ErrNotEnoughArguments = errors.New("NOT ENOUGH ARGUMENTS")
var ErrUnknownCommand = errors.New("UNKNOWN COMMAND")

type Executable interface {
	Execute() (result string)
	Undo()
}

// "SET" command structure
type Set struct {
	Key       string
	Value     string

	old_value string
	existed   bool // Store true if Key already existed is DB
}

// "GET" command structure
type Get struct {
	Key string
}

// "UNSET" command structure
type Unset struct {
	Key       string

	old_value string
	existed   bool
}

// "NUMEQUALTO" command structure
type NumEqualTo struct {
	Value string
}

// Represents group of commands
type Transaction struct {
	commands []Executable
}

//
type NestedTransactions struct {
	transactions []*Transaction
}

// Add command to ast opened transaction
func (self *NestedTransactions) AddCommand(cmd Executable) {
	tr, err := self.getLast()
	if err != nil {
		return
	}
	tr.AddCommand(cmd)
}

// Open new transaction
func (self *NestedTransactions) Begin() (*Transaction) {
	tr := &Transaction{}
	self.transactions = append(self.transactions, tr)
	return tr
}

// Return last opened transaction and delete it from the list
func (self *NestedTransactions) popLast() (*Transaction, error) {
	if len(self.transactions) == 0 {
		return nil, ErrNoTransaction
	}

	last_tr_index := len(self.transactions) - 1
	tr, transactions := self.transactions[last_tr_index], self.transactions[:last_tr_index]
	self.transactions = transactions
	return tr, nil
}

// Commit last opened transaction
func (self *NestedTransactions) Commit() error {
	_, err := self.popLast();
	if err != nil {
		return err
	}
	return nil
}

// Return last opened transaction
func (self *NestedTransactions) getLast() (*Transaction, error) {
	if len(self.transactions) == 0 {
		return nil, ErrNoTransaction
	}

	last_tr_index := len(self.transactions) - 1
	tr := self.transactions[last_tr_index]
	return tr, nil
}

// Rollback last opened transaction
func (self *NestedTransactions) Rollback() error {
	tr, err := self.popLast();
	if err != nil {
		return err
	}
	tr.Rollback()
	return nil
}

// Undo all commands in transaction
func (tr *Transaction) Rollback() {
	for i := len(tr.commands) - 1; i >= 0; i-- {
		cmd := tr.commands[i]
		cmd.Undo()
	}
}

// Add command in transaction
func (tr *Transaction) AddCommand(cmd Executable) {
	tr.commands = append(tr.commands, cmd)
}

// Set value to the key in storage
func (cmd *Set) Execute() (result string) {
	old_value, ok := dbStorage[cmd.Key]
	if ok {
		cmd.old_value = old_value
		cmd.existed = true
	}
	dbStorage[cmd.Key] = cmd.Value
	return
}

// Back key in previous state
func (cmd *Set) Undo() {
	if cmd.existed {
		dbStorage[cmd.Key] = cmd.old_value
	}else {
		delete(dbStorage, cmd.Key)
	}
}

// Delete key from storage
func (cmd *Unset) Execute() (result string) {
	old_value, ok := dbStorage[cmd.Key]
	if ok {
		cmd.existed = true
		cmd.old_value = old_value
		delete(dbStorage, cmd.Key)
	}
	return
}

// Undo deleting key from storage
func (cmd *Unset) Undo() {
	if cmd.existed {
		dbStorage[cmd.Key] = cmd.old_value
	}
}

// Get value by key from storage
func (cmd *Get) Execute() (result string) {
	if value, exists := dbStorage[cmd.Key]; exists {
		return value
	}else {
		return "NULL"
	}
}

// Do nothing, just for implementing Executable interface
func (cmd *Get) Undo() {}


func (cmd *NumEqualTo) Execute() (result string) {
	count := 0
	for _, value := range dbStorage {
		if value == cmd.Value {
			count++
		}
	}
	return strconv.Itoa(count)
}

// Do nothing, just for implementing Executable interface
func (cmd *NumEqualTo) Undo() {}

// Print error if exists
func PrintErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func RunDBCommand(cmd_text string, cmd_args...string) error {
	switch cmd_text {
	case "SET":
		if len(cmd_args) < 2 {
			return ErrNotEnoughArguments
		}

		cmd := &Set{Key: cmd_args[0], Value: cmd_args[1]}
		transactions.AddCommand(cmd)
		cmd.Execute()

	case "GET":
		if len(cmd_args) < 1 {
			return ErrNotEnoughArguments
		}

		cmd := &Get{Key: cmd_args[0]}
		fmt.Println(cmd.Execute())

	case "UNSET":
		if len(cmd_args) < 1 {
			return ErrNotEnoughArguments
		}

		cmd := &Unset{Key: cmd_args[0]}
		transactions.AddCommand(cmd)
		cmd.Execute()

	case "NUMEQUALTO":
		if len(cmd_args) < 1 {
			return ErrNotEnoughArguments
		}

		cmd := &NumEqualTo{Value: cmd_args[0]}
		fmt.Println(cmd.Execute())

	case "BEGIN":
		transactions.Begin()

	case "ROLLBACK":
		err := transactions.Rollback()
		PrintErr(err)

	case "COMMIT":
		err := transactions.Commit()
		PrintErr(err)

	default:
		return ErrUnknownCommand
	}
	return nil
}

func main() {
	for {
		reader := bufio.NewReader(os.Stdin)
		input_text, _ := reader.ReadString('\n')
		input_text = strings.TrimSpace(input_text)
		splited_input := strings.Split(input_text, " ")
		cmd_text, cmd_args := splited_input[0], splited_input[1:]
		cmd_text = strings.ToUpper(strings.TrimSpace(cmd_text))

		if cmd_text == "END" {
			break
		}

		err := RunDBCommand(cmd_text, cmd_args...)
		PrintErr(err)
	}
}