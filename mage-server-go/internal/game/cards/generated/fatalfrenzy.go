package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fatal Frenzy", NewFatalFrenzy)
}

// NewFatalFrenzy creates a Fatal Frenzy
// {2}{R} - INSTANT
func NewFatalFrenzy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fatal Frenzy")
	card.ManaCost = "{2}{R}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeTargetEffect("sacrifice this", source.getControllerId())
	// card.AddAbility(ability0)
	return card, nil
}
