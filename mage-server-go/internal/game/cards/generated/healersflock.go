package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Healers Flock", NewHealersFlock)
}

// NewHealersFlock creates a Healers Flock
// {W}{W}{W} - CREATURE
// Flying, Lifelink
func NewHealersFlock(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Healers Flock")
	card.ManaCost = "{W}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BIRD"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability1)
	return card, nil
}
