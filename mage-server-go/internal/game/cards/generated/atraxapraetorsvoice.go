package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Atraxa Praetors Voice", NewAtraxaPraetorsVoice)
}

// NewAtraxaPraetorsVoice creates a Atraxa Praetors Voice
// {G}{W}{U}{B} - CREATURE
// Flying, Vigilance, Deathtouch, Lifelink
func NewAtraxaPraetorsVoice(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Atraxa Praetors Voice")
	card.ManaCost = "{G}{W}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"PHYREXIAN", "ANGEL", "HORROR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	ability2 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDeathtouch)
	card.AddAbility(ability2)
	ability3 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability3)
	return card, nil
}
