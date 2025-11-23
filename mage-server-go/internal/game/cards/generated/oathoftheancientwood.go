package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Oath Of The Ancient Wood", NewOathOfTheAncientWood)
}

// NewOathOfTheAncientWood creates a Oath Of The Ancient Wood
// {2}{G} - ENCHANTMENT
func NewOathOfTheAncientWood(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Oath Of The Ancient Wood")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
