package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Zevlor Elturel Exile", NewZevlorElturelExile)
}

// NewZevlorElturelExile creates a Zevlor Elturel Exile
// {1}{U}{B}{R} - CREATURE
// Haste
func NewZevlorElturelExile(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Zevlor Elturel Exile")
	card.ManaCost = "{1}{U}{B}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"TIEFLING", "WARRIOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	return card, nil
}
