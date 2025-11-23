package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Silverback Ape", NewSilverbackApe)
}

// NewSilverbackApe creates a Silverback Ape
// {3}{G}{G} - CREATURE
func NewSilverbackApe(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Silverback Ape")
	card.ManaCost = "{3}{G}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"APE"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
