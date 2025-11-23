package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Lavinia Foil To Conspiracy", NewLaviniaFoilToConspiracy)
}

// NewLaviniaFoilToConspiracy creates a Lavinia Foil To Conspiracy
// {1}{W/U}{W/U} - CREATURE
// Vigilance
func NewLaviniaFoilToConspiracy(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lavinia Foil To Conspiracy")
	card.ManaCost = "{1}{W/U}{W/U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "DETECTIVE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	return card, nil
}
