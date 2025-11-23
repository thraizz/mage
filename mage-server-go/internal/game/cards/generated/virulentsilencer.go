package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Virulent Silencer", NewVirulentSilencer)
}

// NewVirulentSilencer creates a Virulent Silencer
// {3} - ARTIFACT CREATURE
func NewVirulentSilencer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Virulent Silencer")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"ROBOT", "ASSASSIN"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
