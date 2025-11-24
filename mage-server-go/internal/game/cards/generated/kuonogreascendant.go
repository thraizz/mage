package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kuon Ogre Ascendant", NewKuonOgreAscendant)
}

// NewKuonOgreAscendant creates a Kuon Ogre Ascendant
// {B}{B}{B} - CREATURE
func NewKuonOgreAscendant(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kuon Ogre Ascendant")
	card.ManaCost = "{B}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"OGRE", "MONK"}
	card.Supertypes = []string{"LEGENDARY", "LEGENDARY"}
	card.Power = "2"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}