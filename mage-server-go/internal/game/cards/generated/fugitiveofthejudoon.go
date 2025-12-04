package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fugitive Of The Judoon", NewFugitiveOfTheJudoon)
}

// NewFugitiveOfTheJudoon creates a Fugitive Of The Judoon
// {4}{G} - ENCHANTMENT
func NewFugitiveOfTheJudoon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fugitive Of The Judoon")
	card.ManaCost = "{4}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"SAGA"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
