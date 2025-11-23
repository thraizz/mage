package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Indomitable", NewTheIndomitable)
}

// NewTheIndomitable creates a The Indomitable
// {2}{U}{U} - ARTIFACT
// Trample
func NewTheIndomitable(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Indomitable")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"VEHICLE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
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
