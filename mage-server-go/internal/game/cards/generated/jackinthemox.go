package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jack In The Mox", NewJackInTheMox)
}

// NewJackInTheMox creates a Jack In The Mox
// {0} - ARTIFACT
func NewJackInTheMox(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jack In The Mox")
	card.ManaCost = "{0}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: SimpleManaAbility
	//   - Effect: JackInTheMoxManaEffect()
	// card.AddAbility(ability0)
	return card, nil
}
