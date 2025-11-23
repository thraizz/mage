package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bond Of Flourishing", NewBondOfFlourishing)
}

// NewBondOfFlourishing creates a Bond Of Flourishing
// {1}{G} - SORCERY
func NewBondOfFlourishing(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bond Of Flourishing")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(3, 1, filter, PutCards.HAND, PutCards.BOTTOM_ANY)
	// card.AddAbility(ability0)
	return card, nil
}
