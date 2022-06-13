package main

import (
	"HW80/hw80Impl"
	"io"
	"os"
)

type Arguments = hw80Impl.Arguments

func Perform(args Arguments, stdOut io.Writer) error {
	return hw80Impl.Perform(args, hw80Impl.GetOSFactory(), stdOut)
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

func parseArgs() Arguments {
	return hw80Impl.NewArgumentsFromCommandLine()
}
