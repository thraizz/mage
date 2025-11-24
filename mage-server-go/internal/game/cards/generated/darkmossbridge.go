package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Darkmoss Bridge", NewDarkmossBridge)
}

// NewDarkmossBridge creates a Darkmoss Bridge
//  - ARTIFACT LAND
// Indestructible
func NewDarkmossBridge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Darkmoss Bridge")
	card.ManaCost = ""
	card.Types = []string{"ARTIFACT", "LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordIndestructible)
	card.AddAbility(ability0)
	ability1 := abilities.BuildSimpleManaAbility(card.ID, "B")
	card.AddAbility(ability1)
	ability2 := abilities.BuildSimpleManaAbility(card.ID, "G")
	card.AddAbility(ability2)
	return card, nil
}