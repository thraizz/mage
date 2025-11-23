package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Palinchron", NewPalinchron)
}

// NewPalinchron creates a Palinchron
// {5}{U}{U} - CREATURE
// Flying
func NewPalinchron(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Palinchron")
	card.ManaCost = "{5}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ILLUSION"}
	card.Power = "4"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
