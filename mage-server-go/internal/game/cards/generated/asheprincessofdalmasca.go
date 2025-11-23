package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ashe Princess Of Dalmasca", NewAshePrincessOfDalmasca)
}

// NewAshePrincessOfDalmasca creates a Ashe Princess Of Dalmasca
// {2}{W} - CREATURE
func NewAshePrincessOfDalmasca(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ashe Princess Of Dalmasca")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "REBEL", "NOBLE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(                 5, 1, StaticFilters.FILTER_CARD_A...)
	// card.AddAbility(ability0)
	return card, nil
}
