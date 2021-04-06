package flags

import (
	"fmt"
	"testing"
)

var opts struct {
	// Slice of bool will append 'true' each time the option
	// is encountered (can be set multiple times, like -vvv)
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`

	// Example of automatic marshalling to desired type (uint)
	Offset uint `long:"offset" description:"Offset"`

	// Example of a required flag
	Name string `short:"n" long:"name" env:"NAME" description:"A name"`

	// Example of a flag restricted to a pre-defined set of strings
	Animal string `long:"animal" choice:"cat" choice:"dog"`

	// Example of a value name
	File string `short:"f" long:"file" description:"A file" value-name:"FILE"`

	// Example of a value name
	Config string `long:"config" short:"c" default:"config.toml" description:"config file" value-name:"FILE"`

	// Example of a pointer
	Ptr *int `short:"p" description:"A pointer to an integer"`

	// Example of a slice of strings
	StringSlice []string `short:"s" description:"A slice of strings"`

	// Example of a slice of pointers
	PtrSlice []*string `long:"ptrslice" description:"A slice of pointers to string"`

	// Example of a map
	IntMap map[string]int `long:"intmap" description:"A map from string to int"`
}

func TestName(t *testing.T) {
	// Make some fake arguments to parse.
	args := []string{
		"-vv",
		"--offset=5",
		//"-n", "Me",
		"--animal", "dog", // anything other than "cat" or "dog" will raise an error
		"-p", "3",
		"-s", "hello",
		"-s", "world",
		//"--ptrslice", "hello",
		//"--ptrslice", "world",
		//"--intmap", "a:1",
		//"--intmap", "b:5",
		//"--config", "config.toml",
		"arg1",
		"arg2",
		"arg3",
	}

	/*_, err := ParseArgs(&opts, args)
	if err != nil {
		panic(err)
	}
	fmt.Println(opts.Offset)*/
	/*p := NewParser(nil, Default)
	ip := NewIniParser(p)
	err = ip.ParseFile("config.ini")
	if err != nil {
		panic(err)
	}*/
	p := NewParser(&opts, Default)
	//os.Setenv("NAME", "oleg")
	if err := NewTomlParser(p).ParseWithConfigOpt("config", args, TomlDefault); err != nil {
		panic(err)
	}

	fmt.Println(opts.Name)
	fmt.Println(opts.IntMap)
}
