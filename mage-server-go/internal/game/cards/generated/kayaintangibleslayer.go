package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kaya Intangible Slayer", NewKayaIntangibleSlayer)
}

// NewKayaIntangibleSlayer creates a Kaya Intangible Slayer
// {3}{W}{W}{B}{B} - PLANESWALKER
// Hexproof
func NewKayaIntangibleSlayer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kaya Intangible Slayer")
	card.ManaCost = "{3}{W}{W}{B}{B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"KAYA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHexproof)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
