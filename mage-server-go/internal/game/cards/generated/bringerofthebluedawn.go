package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Bringer Of The Blue Dawn", NewBringerOfTheBlueDawn)
}

// NewBringerOfTheBlueDawn creates a Bringer Of The Blue Dawn
// {7}{U}{U} - CREATURE
// Trample
func NewBringerOfTheBlueDawn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Bringer Of The Blue Dawn")
	card.ManaCost = "{7}{U}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BRINGER"}
	card.Power = "5"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(2)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
