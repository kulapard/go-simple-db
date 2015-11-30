# go-simple-db

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/kulapard/go-simple-db/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/kulapard/go-simple-db.svg)](https://travis-ci.org/kulapard/go-simple-db)
[![Code Health](https://landscape.io/github/kulapard/go-simple-db/master/landscape.svg?style=flat)](https://landscape.io/github/kulapard/go-simple-db/master)
[![codecov.io](https://codecov.io/github/kulapard/go-simple-db/coverage.svg?branch=master)](https://codecov.io/github/kulapard/go-simple-db?branch=master)
[![Coverage Status](https://coveralls.io/repos/kulapard/go-simple-db/badge.svg?branch=master&service=github)](https://coveralls.io/github/kulapard/go-simple-db?branch=master)

Simple in-memory database written in Go. Same as Redis but much simpler.

**IMPORTANT:** This is not a real database. The only reason I made it is to show my Go programming skills.

> Problem: Simple Database

> Implement an in-memory database similar to Redis. For simplicity your program will receive commands via stdin, and should write appropriate responses to stdout. Client must also support interactive input.

## Requirements
Go (1.5.1 or higher). See installation instruction on official Golang site
([https://golang.org/doc/install](https://golang.org/doc/install)).

## Installation
```
$ go get github.com/kulapard/go-simple-db
$ go install github.com/kulapard/go-simple-db
```
## Running

```
$ go-simple-db
```

## Testing

```
$ go test -v github.com/kulapard/go-simple-db
```

## Commands
### Data Commands

- `SET <name> <value>` – set the variable name to the value value. Neither variable names nor values will contain spaces.
- `GET <name>` – print out the value of the variable name, or `NULL` if that variable is not set.
- `UNSET <name>` – unset the variable name, making it just like that variable was never set.
- `NUMEQUALTO <value>` – print out the number of variables that are currently set to value. If no variables equal that value, print `0`.
- `END` – exit the program.

### Transaction Commands

- `BEGIN` – open a new transaction block. Transaction blocks can be nested; a `BEGIN` can be issued inside of an existing block.
- `ROLLBACK` – undo all of the commands issued in the most recent transaction block, and close the block. Print nothing if successful, or print `NO TRANSACTION` if no transaction is in progress.
- `COMMIT` – close all open transaction blocks, permanently applying the changes made in them. Print nothing if successful, or print `NO TRANSACTION` if no transaction is in progress.

### Examples
Whithout transactions:
```
INPUT          OUTPUT
SET ex 10
GET ex         10
UNSET ex
GET ex         NULL
END
```
```
INPUT          OUTPUT
SET a 10
SET b 10
NUMEQUALTO 10  2
NUMEQUALTO 20  0
SET b 30
NUMEQUALTO 10  1
END
```

With nested transactions:
```
INPUT        OUTPUT
BEGIN
SET a 10     10 
GET a
BEGIN
SET a 20
GET a        20
ROLLBACK
GET a        10
ROLLBACK
GET a        NULL
END
```

```
INPUT        OUTPUT
BEGIN
SET a 30
BEGIN
SET a 40
COMMIT
GET a        40
ROLLBACK     NO TRANSACTION
END
```
```
INPUT        OUTPUT
SET a 50
BEGIN
GET a        50
SET a 60
BEGIN
UNSET a
GET a        NULL
ROLLBACK
GET a        60
COMMIT
GET a        60
END
```
```
INPUT          OUTPUT
SET a 10
BEGIN
NUMEQUALTO 10  1
BEGIN
UNSET a
NUMEQUALTO 10  0
ROLLBACK
NUMEQUALTO 10  1
COMMIT
END
```
