package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Shadowgrange Archfiend", NewShadowgrangeArchfiend)
}

// NewShadowgrangeArchfiend creates a Shadowgrange Archfiend
// {6}{B} - CREATURE
func NewShadowgrangeArchfiend(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shadowgrange Archfiend")
	card.ManaCost = "{6}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DEMON"}
	card.Power = "8"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
