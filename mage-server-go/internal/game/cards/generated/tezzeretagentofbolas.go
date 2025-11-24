package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tezzeret Agent Of Bolas", NewTezzeretAgentOfBolas)
}

// NewTezzeretAgentOfBolas creates a Tezzeret Agent Of Bolas
// {2}{U}{B} - PLANESWALKER
func NewTezzeretAgentOfBolas(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tezzeret Agent Of Bolas")
	card.ManaCost = "{2}{U}{B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"TEZZERET"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 5, 1, StaticFilters.FILTER_CARD_A...)
	// card.AddAbility(ability0)
	return card, nil
}
