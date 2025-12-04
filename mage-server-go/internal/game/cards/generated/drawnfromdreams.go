package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Drawn From Dreams", NewDrawnFromDreams)
}

// NewDrawnFromDreams creates a Drawn From Dreams
// {2}{U}{U} - SORCERY
func NewDrawnFromDreams(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Drawn From Dreams")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(7, 2, PutCards.HAND, PutCards.BOTTOM_RANDOM)
	// card.AddAbility(ability0)
	return card, nil
}
