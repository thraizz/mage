package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Component Pouch", NewComponentPouch)
}

// NewComponentPouch creates a Component Pouch
// {3} - ARTIFACT
func NewComponentPouch(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Component Pouch")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: SimpleManaAbility
	//   - Effect: AddManaOfTwoDifferentColorsEffect()
	// card.AddAbility(ability0)
	return card, nil
}
