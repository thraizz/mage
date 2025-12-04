package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Phyrexian Portal", NewPhyrexianPortal)
}

// NewPhyrexianPortal creates a Phyrexian Portal
// {3} - ARTIFACT
func NewPhyrexianPortal(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Phyrexian Portal")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - PhyrexianPortalEffect()
	// card.AddAbility(ability0)
	return card, nil
}
