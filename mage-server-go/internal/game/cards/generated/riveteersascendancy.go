package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Riveteers Ascendancy", NewRiveteersAscendancy)
}

// NewRiveteersAscendancy creates a Riveteers Ascendancy
// {B}{R}{G} - ENCHANTMENT
func NewRiveteersAscendancy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Riveteers Ascendancy")
	card.ManaCost = "{B}{R}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: SacrificePermanentTriggeredAbility
	//   - Effect: ReturnFromGraveyardToBattlefieldTargetEffect(true)
	// card.AddAbility(ability0)
	return card, nil
}
