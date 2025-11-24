package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Consult The Star Charts", NewConsultTheStarCharts)
}

// NewConsultTheStarCharts creates a Consult The Star Charts
// {1}{U} - INSTANT
func NewConsultTheStarCharts(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Consult The Star Charts")
	card.ManaCost = "{1}{U}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(LandsYouControlCount.instance, 2, PutCards.HAND, P...)
	//   - LookLibraryAndPickControllerEffect(LandsYouControlCount.instance, 1, PutCards.HAND, P...)
	// card.AddAbility(ability0)
	return card, nil
}