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

type dealerAI struct{}

func (ai dealerAI) Bet() int {
	//noop
	return 1
}

func (ai dealerAI) Play(hand []cards.Card, dealer cards.Card) Move {
	dScore := Score(hand...)
	if dScore <= 16 || (dScore == 17 && Soft(hand...)) {
		return MoveHit
	}
	return MoveStand
}

func (ai dealerAI) Results(hand [][]cards.Card, dealer []cards.Card) {
	//noop
}

func HumanAI() AI {
	return humanAI{}
}

type humanAI struct{}

func (ai humanAI) Bet() int {
	return 1
}

func (ai humanAI) Play(hand []cards.Card, dealer cards.Card) Move {
	for {
		fmt.Println("Player:", hand)
		fmt.Println("Dealer:", dealer)
		fmt.Println("What will you do? (h)it, (s)tand")
		var input string
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		default:
			fmt.Println("Invalid option:", input)
		}
	}
}

func (ai humanAI) Results(hand [][]cards.Card, dealer []cards.Card) {
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Player:", hand)
	fmt.Println("Dealer:", dealer)
}
