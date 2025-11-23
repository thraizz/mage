package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vesuvan Doppelganger", NewVesuvanDoppelganger)
}

// NewVesuvanDoppelganger creates a Vesuvan Doppelganger
// {3}{U}{U} - CREATURE
func NewVesuvanDoppelganger(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vesuvan Doppelganger")
	card.ManaCost = "{3}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SHAPESHIFTER"}
	card.Power = "0"
	card.Toughness = "0"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
