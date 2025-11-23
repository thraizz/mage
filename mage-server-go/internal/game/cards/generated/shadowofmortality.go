package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Shadow Of Mortality", NewShadowOfMortality)
}

// NewShadowOfMortality creates a Shadow Of Mortality
// {13}{B}{B} - CREATURE
func NewShadowOfMortality(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shadow Of Mortality")
	card.ManaCost = "{13}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"AVATAR"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
