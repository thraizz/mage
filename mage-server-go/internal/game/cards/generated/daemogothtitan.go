package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Daemogoth Titan", NewDaemogothTitan)
}

// NewDaemogothTitan creates a Daemogoth Titan
// {B/G}{B/G}{B/G}{B/G} - CREATURE
func NewDaemogothTitan(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Daemogoth Titan")
	card.ManaCost = "{B/G}{B/G}{B/G}{B/G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DEMON"}
	card.Power = "11"
	card.Toughness = "10"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement spell ability with unmapped effects
	//   - SacrificeControllerEffect(                 StaticFilters.FILTER_PERMANENT_CR...)
	// card.AddAbility(ability0)
	return card, nil
}