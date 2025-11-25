package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cartographers Survey", NewCartographersSurvey)
}

// NewCartographersSurvey creates a Cartographers Survey
// {3}{G} - SORCERY
func NewCartographersSurvey(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cartographers Survey")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 7, 2, StaticFilters.FILTER_CARD_L...)
	// card.AddAbility(ability0)
	return card, nil
}
