package main

import (
	// "bufio"
	"fmt"
	"os"
	// "strings"
	"shuffler/app"

	term "github.com/nsf/termbox-go"
)

func printWrongCards(wrongCards map[string]int) {
	fmt.Println("Wrong cards: ")
	for k, v := range wrongCards {
		fmt.Printf(" %s: %d times\n", k, v)
	}
}

func main() {
	err := term.Init()
	if err != nil {
			panic(err)
	}

	fmt.Println("Shuffler")
	fmt.Println("------")

	args := os.Args[1:]

	var infinite bool
	if len(args) == 0 {
		infinite = false
	} else {
		if args[0] == "i" {
			fmt.Println(">>> infinite mode <<<")
			infinite = true
		}
	}

	a := app.NewApp()
	c, output, close := a.Run(infinite)

	go func() {
		for {
			select {
			case card := <- output:
				fmt.Println(card.GetName())
			case <- close:
				fmt.Println("closing app")
				printWrongCards(a.GetWrongCards())
				os.Exit(0)
			}
		}
	}()

	for {
			switch ev := term.PollEvent(); ev.Type {
			case term.EventKey:
					switch ev.Key {
					case term.KeyEnter:
						c <- "e"
					case term.KeyArrowRight:
							fmt.Println("perfect")
							c <- "t"
					case term.KeyArrowLeft:
							fmt.Println("missed")
							c <- "f"

					default:
							// we only want to read a single character or one key pressed event
							fmt.Println("ASCII : ", ev.Ch)

					}
			case term.EventError:
					panic(ev.Err)
			}
	}
}

