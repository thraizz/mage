package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Repurposing Bay", NewRepurposingBay)
}

// NewRepurposingBay creates a Repurposing Bay
// {2}{U} - ARTIFACT
func NewRepurposingBay(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Repurposing Bay")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ActivateAsSorceryActivatedAbility
	//   - Effect: RepurposingBayEffect()
	// card.AddAbility(ability0)
	return card, nil
}
