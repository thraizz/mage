package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Cruel Truths", NewCruelTruths)
}

// NewCruelTruths creates a Cruel Truths
// {3}{B} - INSTANT
func NewCruelTruths(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cruel Truths")
	card.ManaCost = "{3}{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewSurveilEffect(1)).
		AddEffect(abilities.NewDrawCardsEffect(2)).
		AddEffect(abilities.NewLoseLifeEffect(2)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
