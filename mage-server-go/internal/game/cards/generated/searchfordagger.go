package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Search For Dagger", NewSearchForDagger)
}

// NewSearchForDagger creates a Search For Dagger
// {1}{W} - ENCHANTMENT
func NewSearchForDagger(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Search For Dagger")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 6, 1, filter, PutCards.HAND, PutC...)
	// card.AddAbility(ability0)
	return card, nil
}
