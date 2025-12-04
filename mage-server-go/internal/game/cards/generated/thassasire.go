package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Thassas Ire", NewThassasIre)
}

// NewThassasIre creates a Thassas Ire
// {U} - ENCHANTMENT
func NewThassasIre(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Thassas Ire")
	card.ManaCost = "{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement activated ability with unmapped effects
	//   - MayTapOrUntapTargetEffect()
	// card.AddAbility(ability0)
	return card, nil
}
