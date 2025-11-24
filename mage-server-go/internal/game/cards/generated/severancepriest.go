package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Severance Priest", NewSeverancePriest)
}

// NewSeverancePriest creates a Severance Priest
// {W}{B}{G} - CREATURE
// Deathtouch
func NewSeverancePriest(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Severance Priest")
	card.ManaCost = "{W}{B}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DJINN", "CLERIC"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDeathtouch)
	card.AddAbility(ability0)
	return card, nil
}
