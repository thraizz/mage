package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gyruda Doom Of Depths", NewGyrudaDoomOfDepths)
}

// NewGyrudaDoomOfDepths creates a Gyruda Doom Of Depths
// {4}{U/B}{U/B} - CREATURE
func NewGyrudaDoomOfDepths(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gyruda Doom Of Depths")
	card.ManaCost = "{4}{U/B}{U/B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DEMON", "KRAKEN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
