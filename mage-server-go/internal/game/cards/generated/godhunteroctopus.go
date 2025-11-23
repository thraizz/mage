package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Godhunter Octopus", NewGodhunterOctopus)
}

// NewGodhunterOctopus creates a Godhunter Octopus
// {5}{U} - CREATURE
func NewGodhunterOctopus(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Godhunter Octopus")
	card.ManaCost = "{5}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"OCTOPUS"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
