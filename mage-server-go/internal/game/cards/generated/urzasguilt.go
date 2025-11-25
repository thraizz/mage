package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Urzas Guilt", NewUrzasGuilt)
}

// NewUrzasGuilt creates a Urzas Guilt
// {2}{U}{B} - SORCERY
func NewUrzasGuilt(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Urzas Guilt")
	card.ManaCost = "{2}{U}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - DiscardEachPlayerEffect(3, false)
	// card.AddAbility(ability0)
	return card, nil
}
