package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jinxed Ring", NewJinxedRing)
}

// NewJinxedRing creates a Jinxed Ring
// {2} - ARTIFACT
func NewJinxedRing(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jinxed Ring")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - TargetPlayerGainControlSourceEffect()
	// card.AddAbility(ability0)
	return card, nil
}
