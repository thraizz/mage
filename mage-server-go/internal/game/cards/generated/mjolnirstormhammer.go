package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mjolnir Storm Hammer", NewMjolnirStormHammer)
}

// NewMjolnirStormHammer creates a Mjolnir Storm Hammer
// {4} - ARTIFACT
func NewMjolnirStormHammer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mjolnir Storm Hammer")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewEquipAbility(card.ID, "{4}", false)
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
