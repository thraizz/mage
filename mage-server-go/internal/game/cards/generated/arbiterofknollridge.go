package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Arbiter Of Knollridge", NewArbiterOfKnollridge)
}

// NewArbiterOfKnollridge creates a Arbiter Of Knollridge
// {6}{W} - CREATURE
// Vigilance
func NewArbiterOfKnollridge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Arbiter Of Knollridge")
	card.ManaCost = "{6}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GIANT", "WIZARD"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	return card, nil
}
