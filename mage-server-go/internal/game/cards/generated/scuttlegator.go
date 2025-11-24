package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Scuttlegator", NewScuttlegator)
}

// NewScuttlegator creates a Scuttlegator
// {4}{G/U}{G/U} - CREATURE
// Defender
func NewScuttlegator(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Scuttlegator")
	card.ManaCost = "{4}{G/U}{G/U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CRAB", "TURTLE", "CROCODILE"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDefender)
	card.AddAbility(ability0)
	return card, nil
}
