package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nettletooth Djinn", NewNettletoothDjinn)
}

// NewNettletoothDjinn creates a Nettletooth Djinn
// {3}{G} - CREATURE
func NewNettletoothDjinn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nettletooth Djinn")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DJINN"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
