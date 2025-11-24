package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Witchs Cottage", NewWitchsCottage)
}

// NewWitchsCottage creates a Witchs Cottage
//   - LAND
func NewWitchsCottage(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Witchs Cottage")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"SWAMP"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "B")
	card.AddAbility(ability0)
	return card, nil
}
