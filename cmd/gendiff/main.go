package main

import (
	"encoding/json"
	"fmt"
	"os"

	"code/internal/parser"

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
			if c.Args().Len() != 2 {
				return fmt.Errorf("error: two file paths are required")
			}

			filepath1 := c.Args().Get(0)
			filepath2 := c.Args().Get(1)
			format := c.String("format")

			// Парсим файлы
			data1, err := parser.ParseFile(filepath1)
			if err != nil {
				return fmt.Errorf("failed to parse %s: %w", filepath1, err)
			}

			data2, err := parser.ParseFile(filepath2)
			if err != nil {
				return fmt.Errorf("failed to parse %s: %w", filepath2, err)
			}

			// Выводим распарсенные данные для отладки
			fmt.Printf("File 1 (%s):\n", filepath1)
			prettyPrint(data1)
			fmt.Printf("\nFile 2 (%s):\n", filepath2)
			prettyPrint(data2)
			fmt.Printf("\nFormat: %s\n", format)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// prettyPrint выводит map в отформатированном JSON-виде
func prettyPrint(data map[string]interface{}) {
	indent, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(indent))
}
