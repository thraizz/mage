package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Estrids Invocation", NewEstridsInvocation)
}

// NewEstridsInvocation creates a Estrids Invocation
// {2}{U} - ENCHANTMENT
func NewEstridsInvocation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Estrids Invocation")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(                 filter, new EstridsInvocationCopy...)
	// card.AddAbility(ability0)
	return card, nil
}