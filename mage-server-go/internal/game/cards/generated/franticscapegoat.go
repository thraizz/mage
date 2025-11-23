package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Frantic Scapegoat", NewFranticScapegoat)
}

// NewFranticScapegoat creates a Frantic Scapegoat
// {R} - CREATURE
// Haste
func NewFranticScapegoat(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Frantic Scapegoat")
	card.ManaCost = "{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOAT"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	return card, nil
}
