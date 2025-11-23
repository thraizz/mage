package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ultimecia Omnipotent", NewUltimeciaOmnipotent)
}

// NewUltimeciaOmnipotent creates a Ultimecia Omnipotent
//   - CREATURE
//
// Menace
func NewUltimeciaOmnipotent(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ultimecia Omnipotent")
	card.ManaCost = ""
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"NIGHTMARE", "WARLOCK"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "7"
	card.Toughness = "7"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordMenace)
	card.AddAbility(ability0)
	return card, nil
}
