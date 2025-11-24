package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Arixmethes Slumbering Isle", NewArixmethesSlumberingIsle)
}

// NewArixmethesSlumberingIsle creates a Arixmethes Slumbering Isle
// {2}{G}{U} - CREATURE
func NewArixmethesSlumberingIsle(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Arixmethes Slumbering Isle")
	card.ManaCost = "{2}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"KRAKEN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "12"
	card.Toughness = "12"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}