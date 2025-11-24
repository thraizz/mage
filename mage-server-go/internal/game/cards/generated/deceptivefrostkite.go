package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Deceptive Frostkite", NewDeceptiveFrostkite)
}

// NewDeceptiveFrostkite creates a Deceptive Frostkite
// {U}{U} - CREATURE
// Flying, Flying
func NewDeceptiveFrostkite(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Deceptive Frostkite")
	card.ManaCost = "{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRAGON"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - CopyPermanentEffect(filter, applier)
	// card.AddAbility(ability2)
	return card, nil
}
