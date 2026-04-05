// Package main предоставляет точку входа для CLI-утилиты gendiff.
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
		// Печатаем ошибку в stderr и завершаемся с кодом 1 — стандартное поведение CLI
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// NewApp создаёт и настраивает CLI-приложение urfave/cli.
// Вынесено в отдельную функцию для упрощения тестирования без вызова os.Exit.
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
			// Валидация: требуется ровно два аргумента — пути к файлам
			if c.Args().Len() != 2 {
				return fmt.Errorf("error: two file paths are required")
			}

			filepath1 := c.Args().Get(0)
			filepath2 := c.Args().Get(1)
			format := c.String("format")

			// Делегируем логику сравнения библиотеке code
			result, err := code.GenDiff(filepath1, filepath2, format)
			if err != nil {
				return err
			}

			fmt.Println(result)
			return nil
		},
	}
}
