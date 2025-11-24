package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gishath Suns Avatar", NewGishathSunsAvatar)
}

// NewGishathSunsAvatar creates a Gishath Suns Avatar
// {5}{R}{G}{W} - CREATURE
// Trample, Vigilance, Haste
func NewGishathSunsAvatar(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gishath Suns Avatar")
	card.ManaCost = "{5}{R}{G}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DINOSAUR", "AVATAR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "7"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	ability2 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHaste)
	card.AddAbility(ability2)
	return card, nil
}