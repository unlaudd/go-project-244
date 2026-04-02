package main

import (
	"fmt"
	"os"

	"code"

	"github.com/urfave/cli/v2"
)

func main() {
	app := NewApp()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func NewApp() *cli.App {
	return &cli.App{
		Name:  "gendiff",
		Usage: "Compares two configuration files and shows a difference.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   "stylish",
				Usage:   "output format",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 2 {
				return fmt.Errorf("error: two file paths are required")
			}

			filepath1 := c.Args().Get(0)
			filepath2 := c.Args().Get(1)
			format := c.String("format")

			result, err := code.GenDiff(filepath1, filepath2, format)
			if err != nil {
				return err
			}

			fmt.Println(result)
			return nil
		},
	}
}
