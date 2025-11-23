package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Domri Chaos Bringer", NewDomriChaosBringer)
}

// NewDomriChaosBringer creates a Domri Chaos Bringer
// {2}{R}{G} - PLANESWALKER
func NewDomriChaosBringer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Domri Chaos Bringer")
	card.ManaCost = "{2}{R}{G}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"DOMRI"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 4, 2, StaticFilters.FILTER_CARD_C...)
	// card.AddAbility(ability0)
	return card, nil
}
