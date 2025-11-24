package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/counters"
)

func init() {
	cards.Register("Declaration Of Naught", NewDeclarationOfNaught)
}

// NewDeclarationOfNaught creates a Declaration Of Naught
// {U}{U} - ENCHANTMENT
func NewDeclarationOfNaught(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Declaration Of Naught")
	card.ManaCost = "{U}{U}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		// TODO: CounterTargetEffect with complex parameters
		Build()
	card.AddAbility(ability0)
	return card, nil
}
