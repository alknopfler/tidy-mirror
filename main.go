package main

import (
	"github.com/TwiN/go-color"
	"github.com/alknopfler/tidy-mirror/cmd"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func init() {

}

func main() {
	command := newCommand()
	if err := command.Execute(); err != nil {
		log.Fatalf(color.InRed("[ERROR] %s"), err.Error())
	}
}

func newCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "t-mirror",
		Short: "Tidy mirror is a cli created just to make it faster using the mirror command",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(1)
		},
	}

	c.AddCommand(cmd.NewMirror())

	return c
}
