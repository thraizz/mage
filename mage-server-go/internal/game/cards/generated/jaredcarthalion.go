package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Jared Carthalion", NewJaredCarthalion)
}

// NewJaredCarthalion creates a Jared Carthalion
// {W}{U}{B}{R}{G} - PLANESWALKER
func NewJaredCarthalion(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jared Carthalion")
	card.ManaCost = "{W}{U}{B}{R}{G}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"JARED"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: LoyaltyAbility
	//   - Effect: JaredCarthalionBoostEffect()
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewCreatureTargetFilter())
	// card.AddAbility(ability0)
	// TODO: Implement triggered ability: LoyaltyAbility
	//   - Effect: JaredCarthalionUltimateEffect()
	// card.AddAbility(ability1)
	token2_0, err := token.GetToken("KavuAllColorToken")
	if err != nil {
		return nil, err
	}
	ability2, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token2_0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability2)
	return card, nil
}
