package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Flayed One", NewFlayedOne)
}

// NewFlayedOne creates a Flayed One
// {2}{B} - ARTIFACT CREATURE
// Lifelink
func NewFlayedOne(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Flayed One")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"NECRON"}
	card.Power = "4"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordLifelink)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewMillCardsControllerEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
