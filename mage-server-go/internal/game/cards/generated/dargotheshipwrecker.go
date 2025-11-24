package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Dargo The Shipwrecker", NewDargoTheShipwrecker)
}

// NewDargoTheShipwrecker creates a Dargo The Shipwrecker
// {6}{R} - CREATURE
// Trample
func NewDargoTheShipwrecker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Dargo The Shipwrecker")
	card.ManaCost = "{6}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GIANT", "PIRATE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "7"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}