package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Soulflayer", NewSoulflayer)
}

// NewSoulflayer creates a Soulflayer
// {4}{B}{B} - CREATURE
// Flying, FirstStrike, DoubleStrike, Deathtouch, Haste, Hexproof, Indestructible, Lifelink, Reach, Trample, Vigilance
func NewSoulflayer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Soulflayer")
	card.ManaCost = "{4}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DEMON"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability1)
	ability2 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDoubleStrike)
	card.AddAbility(ability2)
	ability3 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDeathtouch)
	card.AddAbility(ability3)
	ability4 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability4)
	ability5 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHexproof)
	card.AddAbility(ability5)
	ability6 := abilities.NewKeywordAbility(card.ID, abilities.KeywordIndestructible)
	card.AddAbility(ability6)
	ability7 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability7)
	ability8 := abilities.NewKeywordAbility(card.ID, abilities.KeywordReach)
	card.AddAbility(ability8)
	ability9 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability9)
	ability10 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability10)
	return card, nil
}
