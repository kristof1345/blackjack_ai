package main

import (
	"blackjack"
	"fmt"

	"github.com/kristof1345/cards"
)

type basicAI struct{}

func (ai *basicAI) Bet(shuffled bool) int {
	panic("not implemented") // TODO: Implement
}

func (ai *basicAI) Play(hand []cards.Card, dealer cards.Card) blackjack.Move {
	panic("not implemented") // TODO: Implement
}

func (ai *basicAI) Results(hands [][]cards.Card, dealer []cards.Card) {
	panic("not implemented") // TODO: Implement
}

func main() {
	opts := blackjack.Options{
		Decks:           3,
		Hands:           2,
		BlackjackPayout: 1.5,
	}
	game := blackjack.New(opts)
	winnings := game.Play(blackjack.HumanAI())
	fmt.Println(winnings)
}
