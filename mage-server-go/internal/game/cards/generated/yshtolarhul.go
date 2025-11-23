package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Yshtola Rhul", NewYshtolaRhul)
}

// NewYshtolaRhul creates a Yshtola Rhul
// {4}{U}{U} - CREATURE
func NewYshtolaRhul(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Yshtola Rhul")
	card.ManaCost = "{4}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CAT", "DRUID"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
