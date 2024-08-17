package main

import (
	"fmt"
	"github.com/ccadden/beat-the-box/deck"
)

type Player struct {
	deck     *deck.Deck
	discards *deck.Deck
}

func main() {
	roundCounter := 0

	p1 := newPlayer()
	p2 := newPlayer()

	ok := dealDemCards(p1, p2)

	if !ok {
		panic("problem initializing game")
	}

	for gameCanContinue(p1, p2) && roundCounter < 10000 {
		if p1.deck.Empty() {
			p1.reshuffle()
		}

		if p2.deck.Empty() {
			p2.reshuffle()
		}

		c1, ok := p1.dealCard()

		if !ok {
			panic("Couldn't deal card for p1")
		}

		c2, ok := p2.dealCard()

		if !ok {
			panic("Couldn't deal card for p2")
		}

		switch {
		case c1 > c2:
			p1.discards.AddCards([]int{c1, c2})
		case c2 > c1:
			p2.discards.AddCards([]int{c1, c2})
		default:
			declareWar(p1, p2)
		}

		roundCounter++
	}

	fmt.Printf("P1 total cards: %v\n", p1.totalCards())
	fmt.Printf("P1 hasLost: %v\n", p1.hasLost())
	fmt.Printf("P2 total cards: %v\n", p2.totalCards())
	fmt.Printf("P2 hasLost: %v\n\n", p2.hasLost())
	fmt.Printf("Took %v rounds\n", roundCounter)
}

func gameCanContinue(p1, p2 *Player) bool {
	if p1.totalCards() == 0 {
		return false
	}

	if p2.totalCards() == 0 {
		return false
	}

	return true
}

func dealDemCards(p1, p2 *Player) bool {
	dealDeck := deck.NewDeck()

	for i := range dealDeck.CardsRemaining() {
		cardToDeal, ok := dealDeck.Deal()

		if !ok {
			return false
		}

		if i%2 == 0 {
			p1.deck.AddCard(cardToDeal)
		} else {
			p2.deck.AddCard(cardToDeal)
		}

	}
	return true
}

func newPlayer() *Player {
	p := Player{}
	p.deck = newEmptyDeck()
	p.discards = newEmptyDeck()

	return &p
}

func newEmptyDeck() *deck.Deck {
	return &deck.Deck{}
}

func (p *Player) hasLost() bool {
	return p.deck.Empty() && p.discards.Empty()
}

func (p *Player) reshuffle() {
	for !p.discards.Empty() {
		card, ok := p.discards.Deal()

		if !ok {
			panic("couldn't reshuffle")
		}

		p.deck.AddCard(card)
	}

	p.deck.Shuffle()
}

func (p *Player) totalCards() int {
	return p.deck.CardsRemaining() + p.discards.CardsRemaining()
}

func (p *Player) numCardsToWager() int {
	if p.totalCards() < 4 {
		return 0
	}

	return 3
}

func (p *Player) dealCard() (int, bool) {
	if p.totalCards() < 1 {
		return 0, false
	}

	if p.deck.CardsRemaining() == 0 {
		p.reshuffle()
	}

	deadCard, ok := p.deck.Deal()

	if !ok {
		return 0, false
	}

	return deadCard, true
}

func (p *Player) dealMultiple(numCards int) ([]int, bool) {
	numCardsToDeal := min(numCards, p.totalCards())
	cardsToReturn := []int{}

	for range numCardsToDeal {
		cardToAdd, ok := p.dealCard()
		if !ok {
			return []int{}, false
		}

		cardsToReturn = append(cardsToReturn, cardToAdd)
	}

	return cardsToReturn, true
}

func declareWar(p1, p2 *Player) bool {
	numCards1 := p1.numCardsToWager()
	numCards2 := p2.numCardsToWager()

	if numCards1 == 0 || numCards2 == 0 {
		// Someone is about to lose anyway, not sure we can even get here
		return true
	}

	wageredCards1, ok := p1.dealMultiple(numCards1)

	if !ok {
		panic("problem dealing multiple cards for p1")
	}

	wageredCards2, ok := p1.dealMultiple(numCards2)

	if !ok {
		panic("problem dealing multiple cards for p1")
	}

	if !ok {
		return false
	}

	c1, ok := p1.dealCard()
	c2, ok := p2.dealCard()

	switch {
	case c1 > c2:
		p1.discards.AddCards([]int{c1, c2})
		p1.discards.AddCards(wageredCards1)
		p1.discards.AddCards(wageredCards2)
	case c2 > c1:
		p2.discards.AddCards([]int{c1, c2})
		p2.discards.AddCards(wageredCards1)
		p2.discards.AddCards(wageredCards2)
	default:
		declareWar(p1, p2)
	}

	return true
}

func (p *Player) numWarCardsToWager() int {
	totalCards := p.totalCards()

	if totalCards == 0 {
		return 0
	}

	cardsToWager := min(totalCards, 3)

	if totalCards <= 3 {
		// we need to have a card to actually flip in the war
		cardsToWager -= 1
	}

	return cardsToWager
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
