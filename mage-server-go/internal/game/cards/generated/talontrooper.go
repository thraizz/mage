package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Talon Trooper", NewTalonTrooper)
}

// NewTalonTrooper creates a Talon Trooper
// {1}{W}{U} - CREATURE
// Flying
func NewTalonTrooper(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Talon Trooper")
	card.ManaCost = "{1}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BIRD", "SCOUT"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
