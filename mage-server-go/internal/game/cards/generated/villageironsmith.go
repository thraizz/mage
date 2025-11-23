package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Village Ironsmith", NewVillageIronsmith)
}

// NewVillageIronsmith creates a Village Ironsmith
// {1}{R} - CREATURE
// FirstStrike
func NewVillageIronsmith(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Village Ironsmith")
	card.ManaCost = "{1}{R}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WEREWOLF"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability0)
	return card, nil
}
