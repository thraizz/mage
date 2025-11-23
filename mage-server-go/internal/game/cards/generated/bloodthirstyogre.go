package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bloodthirsty Ogre", NewBloodthirstyOgre)
}

// NewBloodthirstyOgre creates a Bloodthirsty Ogre
// {2}{B} - CREATURE
func NewBloodthirstyOgre(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bloodthirsty Ogre")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"CREATURE"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
