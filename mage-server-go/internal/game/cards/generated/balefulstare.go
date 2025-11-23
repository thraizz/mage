package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Baleful Stare", NewBalefulStare)
}

// NewBalefulStare creates a Baleful Stare
// {2}{U} - SORCERY
func NewBalefulStare(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Baleful Stare")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - RevealHandTargetEffect()
	//
	// Targets:
	//   - abilities.NewOpponentTargetFilter()
	// card.AddAbility(ability0)
	return card, nil
}
