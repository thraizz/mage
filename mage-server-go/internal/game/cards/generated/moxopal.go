package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mox Opal", NewMoxOpal)
}

// NewMoxOpal creates a Mox Opal
// {0} - ARTIFACT
func NewMoxOpal(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mox Opal")
	card.ManaCost = "{0}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ActivateIfConditionManaAbility
	//   - Effect: AddManaOfAnyColorEffect()
	// card.AddAbility(ability0)
	return card, nil
}
