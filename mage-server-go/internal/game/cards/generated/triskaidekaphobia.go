package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Triskaidekaphobia", NewTriskaidekaphobia)
}

// NewTriskaidekaphobia creates a Triskaidekaphobia
// {3}{B} - ENCHANTMENT
func NewTriskaidekaphobia(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Triskaidekaphobia")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: BeginningOfUpkeepTriggeredAbility
	//   - Effect: TriskaidekaphobiaGainLifeEffect()
	// card.AddAbility(ability0)
	return card, nil
}
