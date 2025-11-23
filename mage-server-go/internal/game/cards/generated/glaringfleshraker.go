package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Glaring Fleshraker", NewGlaringFleshraker)
}

// NewGlaringFleshraker creates a Glaring Fleshraker
// {2}{C} - CREATURE
func NewGlaringFleshraker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Glaring Fleshraker")
	card.ManaCost = "{2}{C}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELDRAZI", "DRONE"}
	card.Power = "2"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("EldraziSpawnToken")
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
