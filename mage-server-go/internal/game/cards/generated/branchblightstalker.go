package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Branchblight Stalker", NewBranchblightStalker)
}

// NewBranchblightStalker creates a Branchblight Stalker
// {1}{G} - CREATURE
func NewBranchblightStalker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Branchblight Stalker")
	card.ManaCost = "{1}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "ELF", "SCOUT"}
	card.Power = "3"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
