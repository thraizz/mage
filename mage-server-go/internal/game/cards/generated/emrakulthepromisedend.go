package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Emrakul The Promised End", NewEmrakulThePromisedEnd)
}

// NewEmrakulThePromisedEnd creates a Emrakul The Promised End
// {13} - CREATURE
// Flying, Trample
func NewEmrakulThePromisedEnd(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Emrakul The Promised End")
	card.ManaCost = "{13}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "13"
	card.Toughness = "13"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability1)
	return card, nil
}
