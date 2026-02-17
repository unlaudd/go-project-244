package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
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
			// Проверка наличия двух аргументов
			if c.Args().Len() != 2 {
				return fmt.Errorf("error: two file paths are required")
			}

			filepath1 := c.Args().Get(0)
			filepath2 := c.Args().Get(1)
			format := c.String("format")

			// Вывод информации для отладки
			fmt.Printf("File 1: %s\n", filepath1)
			fmt.Printf("File 2: %s\n", filepath2)
			fmt.Printf("Format: %s\n", format)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
