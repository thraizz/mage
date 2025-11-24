package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Quicksilver Gargantuan", NewQuicksilverGargantuan)
}

// NewQuicksilverGargantuan creates a Quicksilver Gargantuan
// {5}{U}{U} - CREATURE
func NewQuicksilverGargantuan(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Quicksilver Gargantuan")
	card.ManaCost = "{5}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SHAPESHIFTER"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(new QuicksilverGargantuanCopyApplier())
	// card.AddAbility(ability0)
	return card, nil
}
