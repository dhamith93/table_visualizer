package main

type Database interface {
	Connect(user string, password string, host string, database string)
}
