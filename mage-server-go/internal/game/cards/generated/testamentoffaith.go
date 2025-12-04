package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Testament Of Faith", NewTestamentOfFaith)
}

// NewTestamentOfFaith creates a Testament Of Faith
// {W} - ENCHANTMENT
func NewTestamentOfFaith(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Testament Of Faith")
	card.ManaCost = "{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - SetBasePowerToughnessSourceEffect()
	// card.AddAbility(ability0)
	return card, nil
}
