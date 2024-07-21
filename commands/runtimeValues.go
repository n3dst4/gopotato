package commands

import (
	"os/user"

	"github.com/go-playground/validator/v10"
)

var configFilePath string = ""

// config starts out with default values
var config = Config{
	KeepDays:     7,
	JournalsPath: "journals",
	PagesPath:    "pages",
	CarryOverTodos: CarryOverTodos{
		OnlyIncomplete: true,
	},
}

var validate *validator.Validate

var usr *user.User
