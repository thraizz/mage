package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Murderous Redcap", NewMurderousRedcap)
}

// NewMurderousRedcap creates a Murderous Redcap
// {2}{B/R}{B/R} - CREATURE
func NewMurderousRedcap(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Murderous Redcap")
	card.ManaCost = "{2}{B/R}{B/R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOBLIN", "ASSASSIN"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}