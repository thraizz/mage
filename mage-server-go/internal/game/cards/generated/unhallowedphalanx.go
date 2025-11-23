package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Unhallowed Phalanx", NewUnhallowedPhalanx)
}

// NewUnhallowedPhalanx creates a Unhallowed Phalanx
// {4}{B} - CREATURE
func NewUnhallowedPhalanx(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Unhallowed Phalanx")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ZOMBIE", "SOLDIER"}
	card.Power = "1"
	card.Toughness = "13"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
