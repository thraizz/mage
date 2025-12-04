package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Conjurers Closet", NewConjurersCloset)
}

// NewConjurersCloset creates a Conjurers Closet
// {5} - ARTIFACT
func NewConjurersCloset(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Conjurers Closet")
	card.ManaCost = "{5}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BeginningOfEndStepTriggeredAbility
	//   - Effect: ExileThenReturnTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
