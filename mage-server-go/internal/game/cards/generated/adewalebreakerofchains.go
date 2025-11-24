package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Adewale Breaker Of Chains", NewAdewaleBreakerOfChains)
}

// NewAdewaleBreakerOfChains creates a Adewale Breaker Of Chains
// {1}{U}{B} - CREATURE
func NewAdewaleBreakerOfChains(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Adewale Breaker Of Chains")
	card.ManaCost = "{1}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "ASSASSIN", "PIRATE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - RevealLibraryPickControllerEffect(                         6, 1, filter, PutCards.HA...)
	// card.AddAbility(ability0)
	return card, nil
}
