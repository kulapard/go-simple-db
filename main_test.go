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

func TestCommands(t *testing.T) {
	Assert := BuildAssert(t)

	var cmd Executable
	var key, value string

	// SET
	key = RandStringRunes(10)
	value = RandStringRunes(10)

	cmd = &Set{Key: key, Value: value}
	cmd.Execute()

	Assert("SET", dbStorage[key], value)

	// GET
	cmd = &Get{Key: key}
	Assert("GET", cmd.Execute(), value)

	// UNSET
	cmd = &Unset{Key: key}
	cmd.Execute()
	_, ok := dbStorage[key]
	Assert("UNSET", ok, false)

	cmd.Undo()
	Assert("UNSET", dbStorage[key], value)

	// NUMEQUALTO
	count := rand.Intn(10)
	for i := count; i > 0; i-- {
		key = RandStringRunes(10)
		dbStorage[key] = value
	}

	cmd = &NumEqualTo{Value: value}
	Assert("UNSET", cmd.Execute(), strconv.Itoa(count))
}

func TestTransactions(t *testing.T) {
	Assert := BuildAssert(t)

	var cmd Executable
	var key, value1, value2 string

	key = RandStringRunes(10)
	value1, value2 = RandStringRunes(10), RandStringRunes(10)

	// Single transaction
	cmd = &Set{Key: key, Value: value1}
	transactions.Begin()
	transactions.AddCommand(cmd)
	cmd.Execute()

	Assert("Single Tr", dbStorage[key], value1)

	transactions.Rollback()

	_, ok := dbStorage[key]
	Assert("Single Tr", ok, false)

	// Nested transactions
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

}