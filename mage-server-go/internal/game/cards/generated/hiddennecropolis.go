package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Hidden Necropolis", NewHiddenNecropolis)
}

// NewHiddenNecropolis creates a Hidden Necropolis
//   - LAND
func NewHiddenNecropolis(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Hidden Necropolis")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"CAVE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "B")
	card.AddAbility(ability0)
	return card, nil
}
