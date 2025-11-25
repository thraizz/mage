package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rakshasas Bargain", NewRakshasasBargain)
}

// NewRakshasasBargain creates a Rakshasas Bargain
// {2/B}{2/G}{2/U} - INSTANT
func NewRakshasasBargain(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rakshasas Bargain")
	card.ManaCost = "{2/B}{2/G}{2/U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 4, 2, PutCards.HAND, PutCards.GRA...)
	// card.AddAbility(ability0)
	return card, nil
}
