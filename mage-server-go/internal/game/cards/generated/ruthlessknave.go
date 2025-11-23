package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Ruthless Knave", NewRuthlessKnave)
}

// NewRuthlessKnave creates a Ruthless Knave
// {2}{B} - CREATURE
func NewRuthlessKnave(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Ruthless Knave")
	card.ManaCost = "{2}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ORC", "PIRATE"}
	card.Power = "3"
	card.Toughness = "2"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("TreasureToken")
	if err != nil {
		return nil, err
	}
	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, 2)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
