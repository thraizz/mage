package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Underworld Hermit", NewUnderworldHermit)
}

// NewUnderworldHermit creates a Underworld Hermit
// {4}{B}{B} - CREATURE
func NewUnderworldHermit(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Underworld Hermit")
	card.ManaCost = "{4}{B}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "PEASANT"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("SquirrelToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, DevotionCount.B)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
