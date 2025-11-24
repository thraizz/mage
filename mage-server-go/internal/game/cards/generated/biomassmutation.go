package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Biomass Mutation", NewBiomassMutation)
}

// NewBiomassMutation creates a Biomass Mutation
// {X}{G/U}{G/U} - INSTANT
func NewBiomassMutation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Biomass Mutation")
	card.ManaCost = "{X}{G/U}{G/U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}