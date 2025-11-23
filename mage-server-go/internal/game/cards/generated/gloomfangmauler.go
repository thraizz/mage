package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gloomfang Mauler", NewGloomfangMauler)
}

// NewGloomfangMauler creates a Gloomfang Mauler
// {5}{B}{B} - CREATURE
func NewGloomfangMauler(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gloomfang Mauler")
	card.ManaCost = "{5}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"NIGHTMARE"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
