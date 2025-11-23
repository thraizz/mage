package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Iron Bully", NewIronBully)
}

// NewIronBully creates a Iron Bully
// {3} - ARTIFACT CREATURE
// Menace
func NewIronBully(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Iron Bully")
	card.ManaCost = "{3}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"GOLEM"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordMenace)
	card.AddAbility(ability0)
	return card, nil
}
