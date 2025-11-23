package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Spear Of Leonidas", NewTheSpearOfLeonidas)
}

// NewTheSpearOfLeonidas creates a The Spear Of Leonidas
// {2}{R} - ARTIFACT
func NewTheSpearOfLeonidas(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Spear Of Leonidas")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewEquipAbility(card.ID, "{2}", false)
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
