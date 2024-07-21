package config

import (
	"fmt"
	"lumphammer/gopotato/utils"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/go-playground/validator/v10"
	"github.com/kr/pretty"
)

func ParseConfig(txt string) (*Config, error) {
	// config starts out with default values
	var config = Config{
		KeepDays:     7,
		JournalsPath: "journals",
		PagesPath:    "pages",
		CarryOverTodos: CarryOverTodos{
			OnlyIncomplete: true,
		},
	}

	// deserialize config file into &config, overwriting existing values
	_, err := toml.Decode(string(txt), &config)
	if err != nil {
		return nil, err
	}

	config.RootPath = utils.TildeToHomeDir(config.RootPath)

	config.JournalsPath = fmt.Sprintf("%s/%s", config.RootPath, config.JournalsPath)
	config.PagesPath = fmt.Sprintf("%s/%s", config.RootPath, config.PagesPath)

	// validate config
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(config)
	if err != nil {
		return nil, err
	}
	color.Cyan(pretty.Sprint("Config loaded:", config))
	return &config, nil
}
