package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Sundial Of The Infinite", NewSundialOfTheInfinite)
}

// NewSundialOfTheInfinite creates a Sundial Of The Infinite
// {2} - ARTIFACT
func NewSundialOfTheInfinite(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Sundial Of The Infinite")
	card.ManaCost = "{2}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: ActivateIfConditionActivatedAbility
	//   - Effect: EndTurnEffect()
	// card.AddAbility(ability0)
	return card, nil
}
