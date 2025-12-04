package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Whirling Catapult", NewWhirlingCatapult)
}

// NewWhirlingCatapult creates a Whirling Catapult
// {4} - ARTIFACT
func NewWhirlingCatapult(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Whirling Catapult")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - DamageEverythingEffect()
	// card.AddAbility(ability0)
	return card, nil
}
