package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Fin Clade Fugitives", NewFinCladeFugitives)
}

// NewFinCladeFugitives creates a Fin Clade Fugitives
// {5}{G} - CREATURE
func NewFinCladeFugitives(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Fin Clade Fugitives")
	card.ManaCost = "{5}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "SALAMANDER", "ROGUE"}
	card.Power = "7"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
