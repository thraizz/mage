package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Collected Company", NewCollectedCompany)
}

// NewCollectedCompany creates a Collected Company
// {3}{G} - INSTANT
func NewCollectedCompany(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Collected Company")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(6, 2, filter, PutCards.BATTLEFIELD, PutCards.BOTTO...)
	// card.AddAbility(ability0)
	return card, nil
}
