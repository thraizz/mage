package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Overlord Of The Floodpits", NewOverlordOfTheFloodpits)
}

// NewOverlordOfTheFloodpits creates a Overlord Of The Floodpits
// {3}{U}{U} - ENCHANTMENT CREATURE
// Flying
func NewOverlordOfTheFloodpits(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Overlord Of The Floodpits")
	card.ManaCost = "{3}{U}{U}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"AVATAR", "HORROR"}
	card.Power = "5"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
