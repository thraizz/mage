package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Mazzy Truesword Paladin", NewMazzyTrueswordPaladin)
}

// NewMazzyTrueswordPaladin creates a Mazzy Truesword Paladin
// {1}{R}{G}{W} - CREATURE
func NewMazzyTrueswordPaladin(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Mazzy Truesword Paladin")
	card.ManaCost = "{1}{R}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HALFLING", "KNIGHT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
