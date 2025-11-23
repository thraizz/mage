package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hundred Handed One", NewHundredHandedOne)
}

// NewHundredHandedOne creates a Hundred Handed One
// {2}{W}{W} - CREATURE
// Vigilance
func NewHundredHandedOne(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hundred Handed One")
	card.ManaCost = "{2}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GIANT"}
	card.Power = "3"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	return card, nil
}
