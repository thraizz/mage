package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Elvish Aberration", NewElvishAberration)
}

// NewElvishAberration creates a Elvish Aberration
// {5}{G} - CREATURE
func NewElvishAberration(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Elvish Aberration")
	card.ManaCost = "{5}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "MUTANT"}
	card.Power = "4"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
