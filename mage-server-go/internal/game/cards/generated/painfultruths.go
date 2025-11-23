package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Painful Truths", NewPainfulTruths)
}

// NewPainfulTruths creates a Painful Truths
// {2}{B} - SORCERY
func NewPainfulTruths(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Painful Truths")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(ColorsOfManaSpentToCastCount.getInstance())).
		AddEffect(abilities.NewLoseLifeEffect(ColorsOfManaSpentToCastCount.getInstance())).
		AddEffect(abilities.NewDrawCardsEffect(ColorsOfManaSpentToCastCount.getInstance())).
		AddEffect(abilities.NewLoseLifeEffect(ColorsOfManaSpentToCastCount.getInstance())).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
