package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sisters Of Stone Death", NewSistersOfStoneDeath)
}

// NewSistersOfStoneDeath creates a Sisters Of Stone Death
// {4}{B}{B}{G}{G} - CREATURE
func NewSistersOfStoneDeath(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sisters Of Stone Death")
	card.ManaCost = "{4}{B}{B}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GORGON"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "7"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}