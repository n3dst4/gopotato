package cmd

import "github.com/go-playground/validator/v10"

var configFilePath string = ""

// config starts out with default values
var config = Config{
	KeepDays: 7,
	CarryOverTodos: CarryOverTodos{
		OnlyIncomplete: true,
	},
}

var validate *validator.Validate
