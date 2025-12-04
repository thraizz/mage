package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Trackers Instincts", NewTrackersInstincts)
}

// NewTrackersInstincts creates a Trackers Instincts
// {1}{G} - SORCERY
func NewTrackersInstincts(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Trackers Instincts")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - RevealLibraryPickControllerEffect(                 4, 1, StaticFilters.FILTER_CARD_C...)
	// card.AddAbility(ability0)
	return card, nil
}
