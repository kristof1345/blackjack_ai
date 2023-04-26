package blackjack

import (
	"errors"

	"github.com/kristof1345/cards"
)

type state int8

const (
	statePlayerTurn state = iota
	stateDealerTurn
	stateHandOver
)

type Options struct {
	Decks           int
	Hands           int
	BlackjackPayout float64
}

func New(opts Options) Game {
	g := Game{
		state:    statePlayerTurn,
		dealerAI: &dealerAI{},
		balance:  0,
	}

	if opts.Decks == 0 {
		opts.Decks = 3
	}
	if opts.Hands == 0 {
		opts.Hands = 100
	}
	if opts.BlackjackPayout == 0.0 {
		opts.BlackjackPayout = 1.5
	}
	g.nDecks = opts.Decks
	g.nHands = opts.Hands
	g.blackjackPayout = opts.BlackjackPayout

	return g
}

type Game struct {
	deck            []cards.Card
	nDecks          int
	nHands          int
	blackjackPayout float64
	state           state
	player          []cards.Card
	playerBet       int
	balance         int
	dealer          []cards.Card
	dealerAI        AI
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

func bet(g *Game, ai AI, shuffled bool) {
	bet := ai.Bet(shuffled)
	g.playerBet = bet
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
	g.deck = nil
	min := 52 * g.nDecks / 3

	for i := 0; i < g.nHands; i++ {
		shuffled := false
		if len(g.deck) < min {
			g.deck = cards.New(cards.Deck(g.nDecks), cards.Shuffle)
			shuffled = true
		}
		bet(g, ai, shuffled)
		deal(g)
		if Blackjack(g.dealer...) {
			endHand(g, ai)
			continue
		}
		for g.state == statePlayerTurn {
			hand := make([]cards.Card, len(g.player))
			copy(hand, g.player)
			move := ai.Play(hand, g.dealer[0])
			err := move(g)
			switch err {
			case errBust:
				MoveStand(g)
			case nil:
				//noop
			default:
				panic(err)
			}
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

var (
	errBust = errors.New("hand score exceeded 21")
)

type Move func(*Game) error

func MoveHit(g *Game) error {
	hand := g.currentHand()
	var card cards.Card
	card, g.deck = draw(g.deck)
	*hand = append(*hand, card)
	if Score(*hand...) > 21 {
		return errBust
	}
	return nil
}

func MoveDouble(g *Game) error {
	if len(g.player) != 2 {
		return errors.New("Can only double with a hand of 2 cards")
	}
	g.playerBet *= 2

	MoveHit(g)
	return MoveStand(g)
}

func MoveStand(g *Game) error {
	g.state++
	return nil
}

func draw(deck []cards.Card) (cards.Card, []cards.Card) {
	return deck[0], deck[1:]
}

func endHand(g *Game, ai AI) {
	pScore, dScore := Score(g.player...), Score(g.dealer...)
	pBlackjack, dBlackjack := Blackjack(g.player...), Blackjack(g.dealer...)
	winnings := g.playerBet
	switch {
	case pBlackjack && dBlackjack:
		winnings = 0
	case dBlackjack:
		winnings = -winnings
	case pBlackjack:
		winnings = int(float64(winnings) * g.blackjackPayout)
	case pScore > 21:
		winnings = -winnings
	case dScore > 21:
		//win
	case pScore > dScore:
		//win
	case dScore > pScore:
		winnings = -winnings
	case dScore == pScore:
		winnings = 0
	}
	g.balance += winnings
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

func Blackjack(hand ...cards.Card) bool {
	return len(hand) == 2 && Score(hand...) == 21
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
