package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Eyes Of The Watcher", NewEyesOfTheWatcher)
}

// NewEyesOfTheWatcher creates a Eyes Of The Watcher
// {2}{U} - ENCHANTMENT
func NewEyesOfTheWatcher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Eyes Of The Watcher")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DoIfCostPaid(new ScryEffect(2), new ManaCostsImpl<>("{1}"))
	// card.AddAbility(ability0)
	return card, nil
}
