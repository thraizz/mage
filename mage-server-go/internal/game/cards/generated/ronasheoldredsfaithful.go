package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Rona Sheoldreds Faithful", NewRonaSheoldredsFaithful)
}

// NewRonaSheoldredsFaithful creates a Rona Sheoldreds Faithful
// {1}{U}{B}{B} - CREATURE
func NewRonaSheoldredsFaithful(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rona Sheoldreds Faithful")
	card.ManaCost = "{1}{U}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WIZARD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
