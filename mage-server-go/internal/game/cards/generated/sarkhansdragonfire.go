package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sarkhans Dragonfire", NewSarkhansDragonfire)
}

// NewSarkhansDragonfire creates a Sarkhans Dragonfire
// {3}{R}{R} - SORCERY
func NewSarkhansDragonfire(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sarkhans Dragonfire")
	card.ManaCost = "{3}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 5, 1, filter, PutCards.HAND, PutC...)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewAnyTargetFilter())
	// card.AddAbility(ability0)
	return card, nil
}
