package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("You Meet In A Tavern", NewYouMeetInATavern)
}

// NewYouMeetInATavern creates a You Meet In A Tavern
// {2}{G}{G} - SORCERY
func NewYouMeetInATavern(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "You Meet In A Tavern")
	card.ManaCost = "{2}{G}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 5, Integer.MAX_VALUE, StaticFilte...)
	// card.AddAbility(ability0)
	return card, nil
}
