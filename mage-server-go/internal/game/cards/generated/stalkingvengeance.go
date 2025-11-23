package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Stalking Vengeance", NewStalkingVengeance)
}

// NewStalkingVengeance creates a Stalking Vengeance
// {5}{R}{R} - CREATURE
// Haste
func NewStalkingVengeance(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Stalking Vengeance")
	card.ManaCost = "{5}{R}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"AVATAR"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability0)
	return card, nil
}
