package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Azorius Guildgate", NewAzoriusGuildgate)
}

// NewAzoriusGuildgate creates a Azorius Guildgate
//   - LAND
func NewAzoriusGuildgate(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Azorius Guildgate")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"GATE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability1)
	return card, nil
}
