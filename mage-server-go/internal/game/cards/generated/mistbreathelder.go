package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Mistbreath Elder", NewMistbreathElder)
}

// NewMistbreathElder creates a Mistbreath Elder
// {G} - CREATURE
func NewMistbreathElder(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mistbreath Elder")
	card.ManaCost = "{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"FROG", "WARRIOR"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(                 new AddCountersSourceEffect(Count...)
	// card.AddAbility(ability0)
	return card, nil
}
