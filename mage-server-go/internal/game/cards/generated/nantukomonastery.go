package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Nantuko Monastery", NewNantukoMonastery)
}

// NewNantukoMonastery creates a Nantuko Monastery
//   - LAND
//
// FirstStrike
func NewNantukoMonastery(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Nantuko Monastery")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.Subtypes = []string{"INSECT", "MONK"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "C")
	card.AddAbility(ability1)
	return card, nil
}
