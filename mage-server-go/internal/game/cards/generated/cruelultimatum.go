package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cruel Ultimatum", NewCruelUltimatum)
}

// NewCruelUltimatum creates a Cruel Ultimatum
// {U}{U}{B}{B}{B}{R}{R} - SORCERY
func NewCruelUltimatum(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cruel Ultimatum")
	card.ManaCost = "{U}{U}{B}{B}{B}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardTargetEffect(3)
	//
	// Targets:
	//   - abilities.NewOpponentTargetFilter()
	// card.AddAbility(ability0)
	return card, nil
}
