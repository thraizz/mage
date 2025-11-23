package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bojuka Bog", NewBojukaBog)
}

// NewBojukaBog creates a Bojuka Bog
//   - LAND
func NewBojukaBog(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bojuka Bog")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "B")
	card.AddAbility(ability0)
	return card, nil
}
