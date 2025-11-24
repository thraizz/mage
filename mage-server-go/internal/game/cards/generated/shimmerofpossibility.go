package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Shimmer Of Possibility", NewShimmerOfPossibility)
}

// NewShimmerOfPossibility creates a Shimmer Of Possibility
// {1}{U} - SORCERY
func NewShimmerOfPossibility(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Shimmer Of Possibility")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(4, 1, PutCards.HAND, PutCards.BOTTOM_RANDOM)
	// card.AddAbility(ability0)
	return card, nil
}
