package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Exterminatus", NewExterminatus)
}

// NewExterminatus creates a Exterminatus
// {5}{W}{B} - SORCERY
// Indestructible
func NewExterminatus(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Exterminatus")
	card.ManaCost = "{5}{W}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordIndestructible)
	card.AddAbility(ability0)
	return card, nil
}
