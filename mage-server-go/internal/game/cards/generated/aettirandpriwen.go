package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Aettir And Priwen", NewAettirAndPriwen)
}

// NewAettirAndPriwen creates a Aettir And Priwen
// {6} - ARTIFACT
func NewAettirAndPriwen(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Aettir And Priwen")
	card.ManaCost = "{6}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewEquipAbility(card.ID, "{5}", true)
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
