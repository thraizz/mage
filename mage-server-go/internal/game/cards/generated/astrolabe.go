package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Astrolabe", NewAstrolabe)
}

// NewAstrolabe creates a Astrolabe
// {3} - ARTIFACT
func NewAstrolabe(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Astrolabe")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: SimpleManaAbility
	//   - Effect: AddManaOfAnyColorEffect()
	// card.AddAbility(ability0)
	return card, nil
}
