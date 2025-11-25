package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Lukka Bound To Ruin", NewLukkaBoundToRuin)
}

// NewLukkaBoundToRuin creates a Lukka Bound To Ruin
// {2}{R}{R/G/P}{G} - PLANESWALKER
func NewLukkaBoundToRuin(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Lukka Bound To Ruin")
	card.ManaCost = "{2}{R}{R/G/P}{G}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"LUKKA"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: LoyaltyAbility
	//   - Effect: LukkaBoundToRuinManaEffect()
	// card.AddAbility(ability0)
	token1_0, err := token.GetToken("PhyrexianBeastToxicToken")
	if err != nil {
		return nil, err
	}
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token1_0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
