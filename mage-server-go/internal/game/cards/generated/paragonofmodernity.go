package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Paragon Of Modernity", NewParagonOfModernity)
}

// NewParagonOfModernity creates a Paragon Of Modernity
// {4} - ARTIFACT CREATURE
// Flying
func NewParagonOfModernity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Paragon Of Modernity")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ANGEL", "WARRIOR"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
