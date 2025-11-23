package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Edgars Awakening", NewEdgarsAwakening)
}

// NewEdgarsAwakening creates a Edgars Awakening
// {3}{B}{B} - SORCERY
func NewEdgarsAwakening(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Edgars Awakening")
	card.ManaCost = "{3}{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
