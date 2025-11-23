package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mage Slayer", NewMageSlayer)
}

// NewMageSlayer creates a Mage Slayer
// {1}{R}{G} - ARTIFACT
func NewMageSlayer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mage Slayer")
	card.ManaCost = "{1}{R}{G}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewEquipAbility(card.ID, "{3}", false)
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
