package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("March Of Swirling Mist", NewMarchOfSwirlingMist)
}

// NewMarchOfSwirlingMist creates a March Of Swirling Mist
// {X}{U} - INSTANT
func NewMarchOfSwirlingMist(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "March Of Swirling Mist")
	card.ManaCost = "{X}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - PhaseOutTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
