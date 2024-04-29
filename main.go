package main

import "github.com/tywil04/stego/internal/cli"

func main() {
	cmd := cli.NewCommand("Embed", "embed")
	cmd.NewArgument("Input", "The file to embed data in.", true)
	cmd.NewArgument("Output", "The file to output the image with embeded data in.", true)
	cmd.NewFlag("Verbose", "Should we log everything that is happening", "verbose", "v")

	cli.Run()
}
