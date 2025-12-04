package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Doubling Cube", NewDoublingCube)
}

// NewDoublingCube creates a Doubling Cube
// {2} - ARTIFACT
func NewDoublingCube(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Doubling Cube")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: SimpleManaAbility
	//   - Effect: DoublingCubeEffect()
	// card.AddAbility(ability0)
	return card, nil
}
