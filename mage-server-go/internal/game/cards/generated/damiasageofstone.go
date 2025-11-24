package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Damia Sage Of Stone", NewDamiaSageOfStone)
}

// NewDamiaSageOfStone creates a Damia Sage Of Stone
// {4}{B}{G}{U} - CREATURE
// Deathtouch
func NewDamiaSageOfStone(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Damia Sage Of Stone")
	card.ManaCost = "{4}{B}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GORGON", "WIZARD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDeathtouch)
	card.AddAbility(ability0)
	return card, nil
}