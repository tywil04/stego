package cli

import (
	"os"
	"strings"
)

type Argument struct {
	name        string
	description string
	required    bool
}

type Flag struct {
	name          string
	description   string
	flag          string
	shorthandFlag string
}

type Command struct {
	name        string
	description string
	command     string
	arguments   []*Argument
	flags       []*Flag
	subCommands []*Command
	handler     func(*Command)
}

var commands []*Command

func RegisterCommand(newCommands ...*Command) {
	commands = append(commands, newCommands...)
}

func NewCommand(name, description, command string) *Command {
	cmd := &Command{
		name:        name,
		description: description,
		command:     command,
	}
	commands = append(commands, cmd)
	return cmd
}

func (c *Command) NewFlag(name, description, flag, shorthandFlag string) {
	c.flags = append(c.flags, &Flag{
		name:          name,
		flag:          flag,
		shorthandFlag: shorthandFlag,
		description:   description,
	})
}

func (c *Command) NewArgument(name, description string, required bool) {
	c.arguments = append(c.arguments, &Argument{
		name:        name,
		description: description,
		required:    required,
	})
}

func (c *Command) NewSubCommand(name, description, command string) {
	c.subCommands = append(c.subCommands, &Command{
		name:        name,
		description: description,
		command:     command,
	})
}

func (c *Command) Handle(handler func(*Command)) {
	c.handler = handler
}

func Run() {
	var flags []string
	var args = []string

	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "--") {
			// is a flag
			flags = append(flags, arg)
		} else {
			args = append(args, arg)
		}
	}

	var test = func(args []string, command *Command) {
		if command.command == args[0] {
			if len(command.subCommands) == 0 {
				command.handler(command)
			} else {
				for _, subCommand := range command.subCommands {
					test(args[1:], subCommand)
				}
			}
		}
	}

	for _, command := range commands {
		rootCommand := command
		for index, arg := range args {
			if rootCommand.command == arg {
				if len(rootCommand.subCommands) == 0 {
					rootCommand.handler(rootCommand)
				} else {

				}
			}

		}
	}
}

// func main() {
// 	cmd := NewCommand("Embed", "embed")
// 	cmd.NewArgument("Input", "The file to embed data in.", true)
// 	cmd.NewArgument("Output", "The file to output the image with embeded data in.", true)
// 	cmd.NewFlag("Verbose", "Should we log everything that is happening", "verbose", "v")
// }
