package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hearthfire Hobgoblin", NewHearthfireHobgoblin)
}

// NewHearthfireHobgoblin creates a Hearthfire Hobgoblin
// {R/W}{R/W}{R/W} - CREATURE
// DoubleStrike
func NewHearthfireHobgoblin(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hearthfire Hobgoblin")
	card.ManaCost = "{R/W}{R/W}{R/W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOBLIN", "SOLDIER"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDoubleStrike)
	card.AddAbility(ability0)
	return card, nil
}