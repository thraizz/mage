package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Kellan Daring Traveler", NewKellanDaringTraveler)
}

// NewKellanDaringTraveler creates a Kellan Daring Traveler
// {1}{W} - CREATURE
func NewKellanDaringTraveler(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kellan Daring Traveler")
	card.ManaCost = "{1}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "FAERIE", "SCOUT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("MapToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, xValue)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
