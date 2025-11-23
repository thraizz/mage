package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Weatherlight", NewWeatherlight)
}

// NewWeatherlight creates a Weatherlight
// {4} - ARTIFACT
// Flying
func NewWeatherlight(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Weatherlight")
	card.ManaCost = "{4}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"VEHICLE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(5, 1, filter, PutCards.HAND, PutCards.BOTTOM_RANDO...)
	// card.AddAbility(ability1)
	return card, nil
}
