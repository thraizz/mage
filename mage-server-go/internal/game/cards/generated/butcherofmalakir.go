package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Butcher Of Malakir", NewButcherOfMalakir)
}

// NewButcherOfMalakir creates a Butcher Of Malakir
// {5}{B}{B} - CREATURE
// Flying
func NewButcherOfMalakir(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Butcher Of Malakir")
	card.ManaCost = "{5}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Power = "5"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
