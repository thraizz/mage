package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fang Fearless L Cie", NewFangFearlessLCie)
}

// NewFangFearlessLCie creates a Fang Fearless L Cie
// {2}{B} - CREATURE
func NewFangFearlessLCie(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fang Fearless L Cie")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WARRIOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
