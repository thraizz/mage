package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Quartzwood Crasher", NewQuartzwoodCrasher)
}

// NewQuartzwoodCrasher creates a Quartzwood Crasher
// {2}{R}{R}{G} - CREATURE
// Trample
func NewQuartzwoodCrasher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Quartzwood Crasher")
	card.ManaCost = "{2}{R}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DINOSAUR", "BEAST"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	return card, nil
}
