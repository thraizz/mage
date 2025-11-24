package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Thieves Fortune", NewThievesFortune)
}

// NewThievesFortune creates a Thieves Fortune
// {2}{U} - KINDRED INSTANT
func NewThievesFortune(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Thieves Fortune")
	card.ManaCost = "{2}{U}"
	card.Types = []string{"KINDRED", "INSTANT"}
	card.Subtypes = []string{"ROGUE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(4, 1, PutCards.HAND, PutCards.BOTTOM_ANY)
	// card.AddAbility(ability0)
	return card, nil
}