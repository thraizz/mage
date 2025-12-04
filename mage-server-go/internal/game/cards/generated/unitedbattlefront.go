package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("United Battlefront", NewUnitedBattlefront)
}

// NewUnitedBattlefront creates a United Battlefront
// {3}{W} - SORCERY
func NewUnitedBattlefront(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "United Battlefront")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 7, 2, filter, PutCards.BATTLEFIEL...)
	// card.AddAbility(ability0)
	return card, nil
}
