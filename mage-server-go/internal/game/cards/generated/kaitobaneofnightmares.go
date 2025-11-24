package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kaito Bane Of Nightmares", NewKaitoBaneOfNightmares)
}

// NewKaitoBaneOfNightmares creates a Kaito Bane Of Nightmares
// {2}{U}{B} - PLANESWALKER
// Hexproof
func NewKaitoBaneOfNightmares(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kaito Bane Of Nightmares")
	card.ManaCost = "{2}{U}{B}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"KAITO"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordHexproof)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewSurveilEffect(1)).
		// TODO: TapTargetEffect with complex parameters
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
