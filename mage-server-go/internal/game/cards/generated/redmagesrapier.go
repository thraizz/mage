package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Red Mages Rapier", NewRedMagesRapier)
}

// NewRedMagesRapier creates a Red Mages Rapier
// {1}{R} - ARTIFACT
func NewRedMagesRapier(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Red Mages Rapier")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewEquipAbility(card.ID, "{3}", true)
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
