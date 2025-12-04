package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Storm The Festival", NewStormTheFestival)
}

// NewStormTheFestival creates a Storm The Festival
// {3}{G}{G}{G} - SORCERY
func NewStormTheFestival(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Storm The Festival")
	card.ManaCost = "{3}{G}{G}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(5, 2, filter, PutCards.BATTLEFIELD, PutCards.BOTTO...)
	// card.AddAbility(ability0)
	return card, nil
}
