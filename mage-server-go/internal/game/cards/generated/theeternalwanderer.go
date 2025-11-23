package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("The Eternal Wanderer", NewTheEternalWanderer)
}

// NewTheEternalWanderer creates a The Eternal Wanderer
// {4}{W}{W} - PLANESWALKER
func NewTheEternalWanderer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Eternal Wanderer")
	card.ManaCost = "{4}{W}{W}"
	card.Types = []string{"PLANESWALKER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("DoublestrikeSamuraiToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
