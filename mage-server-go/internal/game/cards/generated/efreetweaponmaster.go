package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Efreet Weaponmaster", NewEfreetWeaponmaster)
}

// NewEfreetWeaponmaster creates a Efreet Weaponmaster
// {3}{U}{R}{W} - CREATURE
// FirstStrike
func NewEfreetWeaponmaster(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Efreet Weaponmaster")
	card.ManaCost = "{3}{U}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"EFREET", "MONK"}
	card.Power = "4"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability0)
	return card, nil
}
