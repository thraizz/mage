package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Trailblazers Torch", NewTrailblazersTorch)
}

// NewTrailblazersTorch creates a Trailblazers Torch
// {4} - ARTIFACT
func NewTrailblazersTorch(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Trailblazers Torch")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewEquipAbility(card.ID, "{1}", true)
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
