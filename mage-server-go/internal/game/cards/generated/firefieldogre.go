package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fire Field Ogre", NewFireFieldOgre)
}

// NewFireFieldOgre creates a Fire Field Ogre
// {1}{U}{B}{R} - CREATURE
// FirstStrike
func NewFireFieldOgre(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fire Field Ogre")
	card.ManaCost = "{1}{U}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"OGRE", "MUTANT"}
	card.Power = "4"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability0)
	return card, nil
}
