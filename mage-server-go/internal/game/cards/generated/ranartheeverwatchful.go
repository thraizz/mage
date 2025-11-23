package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Ranar The Ever Watchful", NewRanarTheEverWatchful)
}

// NewRanarTheEverWatchful creates a Ranar The Ever Watchful
// {2}{W}{U} - CREATURE
// Flying, Vigilance
func NewRanarTheEverWatchful(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ranar The Ever Watchful")
	card.ManaCost = "{2}{W}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT", "WARRIOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	return card, nil
}
