package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Pillage The Bog", NewPillageTheBog)
}

// NewPillageTheBog creates a Pillage The Bog
// {B}{G} - SORCERY
func NewPillageTheBog(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Pillage The Bog")
	card.ManaCost = "{B}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(xValue, 1, PutCards.HAND, PutCards.BOTTOM_RANDOM)
	// card.AddAbility(ability0)
	return card, nil
}
