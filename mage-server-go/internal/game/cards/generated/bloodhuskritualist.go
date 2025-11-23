package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bloodhusk Ritualist", NewBloodhuskRitualist)
}

// NewBloodhuskRitualist creates a Bloodhusk Ritualist
// {2}{B} - CREATURE
func NewBloodhuskRitualist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bloodhusk Ritualist")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"CREATURE"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
