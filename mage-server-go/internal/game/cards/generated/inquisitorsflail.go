package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Inquisitors Flail", NewInquisitorsFlail)
}

// NewInquisitorsFlail creates a Inquisitors Flail
// {2} - ARTIFACT
func NewInquisitorsFlail(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Inquisitors Flail")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewEquipAbility(card.ID, "{2}", false)
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
