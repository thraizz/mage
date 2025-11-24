package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Flubs The Fool", NewFlubsTheFool)
}

// NewFlubsTheFool creates a Flubs The Fool
// {G}{U}{R} - CREATURE
func NewFlubsTheFool(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Flubs The Fool")
	card.ManaCost = "{G}{U}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FROG", "SCOUT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "0"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardControllerEffect(1)
	// card.AddAbility(ability0)
	return card, nil
}
