package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Comet Stellar Pup", NewCometStellarPup)
}

// NewCometStellarPup creates a Comet Stellar Pup
// {2}{R}{W} - PLANESWALKER
func NewCometStellarPup(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Comet Stellar Pup")
	card.ManaCost = "{2}{R}{W}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"COMET"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
