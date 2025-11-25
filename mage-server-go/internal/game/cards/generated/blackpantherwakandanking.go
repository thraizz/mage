package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Black Panther Wakandan King", NewBlackPantherWakandanKing)
}

// NewBlackPantherWakandanKing creates a Black Panther Wakandan King
// {G}{W} - CREATURE
// FirstStrike
func NewBlackPantherWakandanKing(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Black Panther Wakandan King")
	card.ManaCost = "{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "NOBLE", "HERO"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability0)
	// TODO: Implement activated ability with unmapped effects
	//   - BlackPantherWakandanKingEffect()
	//
	// Costs:
	//   - AddManaCost("{3}")
	// card.AddAbility(ability1)
	return card, nil
}
