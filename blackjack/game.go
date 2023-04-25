package blackjack

import (
	"fmt"

	"github.com/kristof1345/cards"
)

type state int8

const (
	statePlayerTurn state = iota
	stateDealerTurn
	stateHandOver
)

func New() Game {
	return Game{
		state:    statePlayerTurn,
		dealerAI: &dealerAI{},
		balance:  0,
	}
}

type Game struct {
	deck     []cards.Card
	state    state
	player   []cards.Card
	dealer   []cards.Card
	dealerAI AI
	balance  int
}

func (g *Game) currentHand() *[]cards.Card {
	switch g.state {
	case statePlayerTurn:
		return &g.player
	case stateDealerTurn:
		return &g.dealer
	default:
		panic("There isnt a turn for a player")
	}
}

func deal(g *Game) {
	g.player = make([]cards.Card, 0, 5)
	g.dealer = make([]cards.Card, 0, 5)
	var card cards.Card
	for i := 0; i < 2; i++ {
		card, g.deck = draw(g.deck)
		g.player = append(g.player, card)
		card, g.deck = draw(g.deck)
		g.dealer = append(g.dealer, card)
	}
	g.state = statePlayerTurn
}

func (g *Game) Play(ai AI) int {
	g.deck = cards.New(cards.Deck(3), cards.Shuffle)
	for i := 0; i < 2; i++ {
		deal(g)

		for g.state == statePlayerTurn {
			hand := make([]cards.Card, len(g.player))
			copy(hand, g.player)
			move := ai.Play(hand, g.dealer[0])
			move(g)
		}
		for g.state == stateDealerTurn {
			hand := make([]cards.Card, len(g.dealer))
			copy(hand, g.dealer)
			move := g.dealerAI.Play(hand, g.dealer[0])
			move(g)
		}
		endHand(g, ai)
	}
	return g.balance
}

type Move func(*Game)

func MoveHit(g *Game) {
	hand := g.currentHand()
	var card cards.Card
	card, g.deck = draw(g.deck)
	*hand = append(*hand, card)
	if Score(*hand...) > 21 {
		MoveStand(g)
	}
}

func MoveStand(g *Game) {
	g.state++
}

func draw(deck []cards.Card) (cards.Card, []cards.Card) {
	return deck[0], deck[1:]
}

func endHand(g *Game, ai AI) {
	pScore, dScore := Score(g.player...), Score(g.dealer...)
	switch {
	case pScore > 21:
		fmt.Println("You busted")
		g.balance--
	case dScore > 21:
		fmt.Println("Dealer busted")
		g.balance++
	case pScore > dScore:
		fmt.Println("You win!")
		g.balance++
	case dScore > pScore:
		fmt.Println("You lose!")
		g.balance--
	case dScore == pScore:
		fmt.Println("Draw")
	}
	fmt.Println()
	ai.Results([][]cards.Card{g.player}, g.dealer)
	g.player = nil
	g.dealer = nil
}

func Score(hand ...cards.Card) int {
	minScore := minScore(hand...)
	if minScore > 11 {
		return minScore
	}
	for _, c := range hand {
		if c.Rank == cards.Ace {
			return minScore + 10
		}
	}
	return minScore
}

func Soft(hand ...cards.Card) bool {
	minScore := minScore(hand...)
	score := Score(hand...)
	return minScore != score
}

func minScore(hand ...cards.Card) int {
	score := 0
	for _, c := range hand {
		score += min(int(c.Rank), 10)
	}
	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
