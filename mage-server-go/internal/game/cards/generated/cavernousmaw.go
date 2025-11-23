package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cavernous Maw", NewCavernousMaw)
}

// NewCavernousMaw creates a Cavernous Maw
//   - LAND
func NewCavernousMaw(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cavernous Maw")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"CAVE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability0)
	return card, nil
}
