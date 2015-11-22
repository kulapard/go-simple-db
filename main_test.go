package main

import (
	"fmt"
	"testing"
	"math/rand"
	"time"
	"strconv"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Returns random string specified length
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Helper to use in tests
func BuildAssert(t *testing.T) func(name string, actual, expected interface{}) {
	return func(prefix string, actual, expected interface{}) {
		if actual != expected {
			t.Error(fmt.Sprintf("%s: %s != %s", prefix, expected, actual))
		}
	}
}

func cleanStorage()  {
	dbStorage = map[string]string{}
}

func TestCommands(t *testing.T) {
	Assert := BuildAssert(t)

	var cmd Executable
	var key, value string
	var ok bool

	// SET
	cleanStorage()
	key = RandStringRunes(10)
	value = RandStringRunes(10)

	cmd = &Set{Key: key, Value: value}
	cmd.Execute()

	Assert("SET", dbStorage[key], value)

	cmd.Undo()
	_, ok = dbStorage[key]
	Assert("SET", ok, false)

	// GET
	cleanStorage()
	dbStorage[key] = value

	cmd = &Get{Key: key}
	Assert("GET", cmd.Execute(), value)

	cmd.Undo() // just for coverage

	cleanStorage()
	key = RandStringRunes(10)
	cmd = &Get{Key: key}
	Assert("GET", cmd.Execute(), "NULL")

	// UNSET
	cleanStorage()
	dbStorage[key] = value

	cmd = &Unset{Key: key}
	cmd.Execute()
	_, ok = dbStorage[key]
	Assert("UNSET", ok, false)

	cmd.Undo()
	Assert("UNSET", dbStorage[key], value)

	// NUMEQUALTO
	cleanStorage()

	count := rand.Intn(10)
	for i := 0; i < count; i++ {
		key = RandStringRunes(10)
		dbStorage[key] = value
	}

	cmd = &NumEqualTo{Value: value}
	Assert("NUMEQUALTO", cmd.Execute(), strconv.Itoa(count))
}

func TestTransactions(t *testing.T) {
	Assert := BuildAssert(t)

	var cmd Executable
	var key, value1, value2 string

	key = RandStringRunes(10)
	value1, value2 = RandStringRunes(10), RandStringRunes(10)

	// Single transaction
	cleanStorage()
	cmd = &Set{Key: key, Value: value1}
	transactions.Begin()
	transactions.AddCommand(cmd)
	cmd.Execute()

	Assert("Single Tr", dbStorage[key], value1)

	transactions.Rollback()

	_, ok := dbStorage[key]
	Assert("Single Tr", ok, false)

	// Nested transactions
	cleanStorage()
	cmd = &Set{Key: key, Value: value1}
	transactions.Begin()
	transactions.AddCommand(cmd)
	cmd.Execute()

	Assert("Nested Tr", dbStorage[key], value1)

	cmd = &Set{Key: key, Value: value2}
	transactions.Begin()
	transactions.AddCommand(cmd)
	cmd.Execute()

	Assert("Nested Tr", dbStorage[key], value2)

	transactions.Rollback()

	Assert("Nested Tr", dbStorage[key], value1)

	transactions.Rollback()

	_, ok = dbStorage[key]
	Assert("Nested Tr", ok, false)

	// Begin/Commit
	cleanStorage()
	transactions.Begin()
	transactions.Begin()
	transactions.Begin()
	Assert("Nested Tr", len(transactions.transactions), 3)
	transactions.Commit()
	Assert("Nested Tr", len(transactions.transactions), 0)

}

func TestRunDBCommand(t *testing.T) {
	Assert := BuildAssert(t)

	var key, value string

	cleanStorage()
	key = RandStringRunes(10)
	value = RandStringRunes(10)

	// Good
	RunDBCommand("SET", key, value)
	RunDBCommand("GET", key)
	RunDBCommand("UNSET", key)
	RunDBCommand("NUMEQUALTO", key)
	RunDBCommand("BEGIN")
	RunDBCommand("COMMIT")
	RunDBCommand("ROLLBACK")

	// Bad
	Assert("Bad", RunDBCommand("SOMEBADCOMMAND"), ErrUnknownCommand)
	Assert("Bad", RunDBCommand("SET", key), ErrNotEnoughArguments)
	Assert("Bad", RunDBCommand("GET"), ErrNotEnoughArguments)
	Assert("Bad", RunDBCommand("UNSET"), ErrNotEnoughArguments)
	Assert("Bad", RunDBCommand("NUMEQUALTO"), ErrNotEnoughArguments)
}