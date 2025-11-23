package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cathedral Membrane", NewCathedralMembrane)
}

// NewCathedralMembrane creates a Cathedral Membrane
// {1}{W/P} - ARTIFACT CREATURE
// Defender
func NewCathedralMembrane(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cathedral Membrane")
	card.ManaCost = "{1}{W/P}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "WALL"}
	card.Power = "0"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDefender)
	card.AddAbility(ability0)
	return card, nil
}
