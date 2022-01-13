package app

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

var (
	// notes      = []string{"A", "E", "D", "G", "B"}
	notes      = []string{"A"}
	stringNums = []int{1, 2, 3, 4, 5, 6}
)

type Card struct {
	Note      string
	StringNum int
}

func NewCard(note string, stringNum int) Card {
	return Card{note, stringNum}
}

func (c Card) GetName() string {
	return fmt.Sprintf("%s-%d", c.Note, c.StringNum)
}

type App struct {
	wrongCards map[string]int
	cards      Cards
	input      chan string
	output     chan Card
	close      chan bool
}

func NewCards() Cards {
	var cards Cards
	for i := range notes {
		for j := range stringNums {
			cards = append(cards, NewCard(notes[i], stringNums[j]))
		}
	}
	shuffleCards(cards)

	return cards
}

func NewApp() App {
	var app App

	app.cards = NewCards()
	app.wrongCards = make(map[string]int)

	return app
}

func shuffleCards(cards []Card) {
	rand.Seed(time.Now().UnixNano())
	for i := len(cards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
}

type Cards []Card

func (cs *Cards) Take() Card {
	card := (*cs)[0]
	(*cs) = (*cs)[1:]

	return card
}

func (a *App) TakeCard() (bool, Card) {
	if len(a.cards) == 0 {
		return false, Card{}
	}
	return true, a.cards.Take()
}

func (a App) GetWrongCards() map[string]int {
	return a.wrongCards
}

func (a *App) Run(infinite bool) (chan string, chan Card, chan bool) {
	input := make(chan string)
	output := make(chan Card)
	close := make(chan bool)

	a.input = input
	a.output = output

	go func() {
		for {
			ok, curCard := a.TakeCard()
			if !ok {
				if infinite {
					a.cards = NewCards()
					_, curCard = a.TakeCard()
				} else {
					fmt.Println("out of cards")
					os.Exit(0)
				}
			}

			output <- curCard

			char := <-input
			switch char {
			case "t":
				continue
			case "f":
				a.wrongCards[curCard.GetName()] += 1
			case "e":
				close <- true
			}
		}
	}()

	return input, output, close
}

func (a *App) End() {
	close(a.input)
	close(a.output)
	close(a.close)
}

func (a *App) GetDeck() Cards {
	return a.cards
}
