package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jinxed Idol", NewJinxedIdol)
}

// NewJinxedIdol creates a Jinxed Idol
// {2} - ARTIFACT
func NewJinxedIdol(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jinxed Idol")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - TargetPlayerGainControlSourceEffect()
	// card.AddAbility(ability0)
	return card, nil
}
