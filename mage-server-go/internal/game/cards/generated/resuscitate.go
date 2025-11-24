package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Resuscitate", NewResuscitate)
}

// NewResuscitate creates a Resuscitate
// {1}{G} - INSTANT
func NewResuscitate(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Resuscitate")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - RegenerateSourceEffect()
	//
	// Costs:
	//   - AddManaCost("{1}")
	// card.AddAbility(ability0)
	return card, nil
}