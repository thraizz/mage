package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Vivien Natures Avenger", NewVivienNaturesAvenger)
}

// NewVivienNaturesAvenger creates a Vivien Natures Avenger
// {4}{G}{G} - PLANESWALKER
func NewVivienNaturesAvenger(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Vivien Natures Avenger")
	card.ManaCost = "{4}{G}{G}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"VIVIEN"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
