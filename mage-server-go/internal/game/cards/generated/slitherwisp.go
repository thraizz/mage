package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Slitherwisp", NewSlitherwisp)
}

// NewSlitherwisp creates a Slitherwisp
// {U}{B}{B} - CREATURE
// Flash
func NewSlitherwisp(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Slitherwisp")
	card.ManaCost = "{U}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELEMENTAL", "NIGHTMARE"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlash)
	card.AddAbility(ability0)
	return card, nil
}
