package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Guidestone Compass", NewGuidestoneCompass)
}

// NewGuidestoneCompass creates a Guidestone Compass
//   - ARTIFACT
func NewGuidestoneCompass(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Guidestone Compass")
	card.ManaCost = ""
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ActivateAsSorceryActivatedAbility
	//   - Effect: ExploreTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
