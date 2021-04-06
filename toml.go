package flags

import (
	"fmt"
	"reflect"

	"github.com/pelletier/go-toml"
)

// TomlOptions for writing
type TomlOptions uint

const (
	// TomlNone indicates no options.
	TomlNone TomlOptions = 0

	// TomlErrOnDefault indicates that default values should be written.
	TomlErrOnDefault = 1 << iota

	// TomlDefault provides a default set of options.
	TomlDefault = TomlErrOnDefault
)

type TomlParser struct {
	parser *Parser
}

func NewTomlParser(p *Parser) *TomlParser {
	return &TomlParser{
		parser: p,
	}
}

func (p *TomlParser) Parse(filename string, args []string, options TomlOptions) error {
	if _, err := p.parser.ParseArgs(args); err != nil {
		return err
	}

	tree, err := toml.LoadFile(filename)
	if err != nil {
		return err
	}

	var optErr error
	p.parser.eachOption(func(command *Command, group *Group, option *Option) {
		if option.IsSet() {
			return
		}
		option.clearReferenceBeforeSet = true
		opName := option.LongNameWithNamespace()
		if opName == "" {
			return
		}

		value := tree.Get(option.LongName)
		if value == nil {
			return
		}

		if subtree, ok := value.(*toml.Tree); ok {
			value = subtree.ToMap()
			fmt.Println(convertToString(reflect.ValueOf(value), option.tag))
		}

		var s string
		s, optErr = convertToString(reflect.ValueOf(value), option.tag)
		if optErr != nil {
			return
		}

		optErr = option.Set(&s)
	})

	return optErr
}

func (p *TomlParser) ParseWithConfigOpt(opName string, args []string, options TomlOptions) error {
	option := p.parser.FindOptionByLongName(opName)
	if option == nil {
		return fmt.Errorf("option %s is not defined", opName)
	}
	if option.Field().Type.Kind() != reflect.String {
		return fmt.Errorf("option %s must be of type string", opName)
	}
	tmpParser := NewParser(nil, p.parser.Options|IgnoreUnknown)
	tmpParser.Group.options = append(tmpParser.Group.options, option)

	_, err := tmpParser.ParseArgs(args)
	if err != nil {
		return err
	}
	fmt.Println(option.Value())

	return p.Parse(option.Value().(string), args, options)
}
