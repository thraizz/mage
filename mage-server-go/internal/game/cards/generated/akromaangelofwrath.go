package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Akroma Angel Of Wrath", NewAkromaAngelOfWrath)
}

// NewAkromaAngelOfWrath creates a Akroma Angel Of Wrath
// {5}{W}{W}{W} - CREATURE
// Flying, FirstStrike, Vigilance, Trample, Haste
func NewAkromaAngelOfWrath(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Akroma Angel Of Wrath")
	card.ManaCost = "{5}{W}{W}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ANGEL"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike)
	card.AddAbility(ability1)
	ability2 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability2)
	ability3 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability3)
	ability4 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability4)
	return card, nil
}
