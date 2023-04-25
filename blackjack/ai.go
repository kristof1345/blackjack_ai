package blackjack

import (
	"fmt"

	"github.com/kristof1345/cards"
)

type AI interface {
	Bet() int
	Play(hand []cards.Card, dealer cards.Card) Move
	Results(hand [][]cards.Card, dealer []cards.Card)
}

type HumanAI struct{}

func (ai *HumanAI) Bet() int {
	return 1
}

func (ai *HumanAI) Play(hand []cards.Card, dealer cards.Card) Move {
	for {
		fmt.Println("Player:", hand)
		fmt.Println("Dealer:", dealer)
		fmt.Println("What will you do? (h)it, (s)tand")
		var input string
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			return Hit
		case "s":
			return Stand
		default:
			fmt.Println("Invalid option:", input)
		}
	}
}

func (ai *HumanAI) Results(hand [][]cards.Card, dealer []cards.Card) {
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Player:", hand)
	fmt.Println("Dealer:", dealer)
}

type Move func(GameState) GameState

type GameState struct{}

func Hit(gs GameState) GameState {
	return gs
}

func Stand(gs GameState) GameState {
	return gs
}
