package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Glacial Revelation", NewGlacialRevelation)
}

// NewGlacialRevelation creates a Glacial Revelation
// {2}{G} - SORCERY
func NewGlacialRevelation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Glacial Revelation")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - RevealLibraryPickControllerEffect(                 6, Integer.MAX_VALUE, filter, Put...)
	// card.AddAbility(ability0)
	return card, nil
}
