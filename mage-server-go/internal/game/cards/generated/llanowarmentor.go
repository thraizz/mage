package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Llanowar Mentor", NewLlanowarMentor)
}

// NewLlanowarMentor creates a Llanowar Mentor
// {G} - CREATURE
func NewLlanowarMentor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Llanowar Mentor")
	card.ManaCost = "{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"ELF", "SPELLSHAPER"}
	card.Power = "1"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("LlanowarElvesToken")
	if err != nil {
		return nil, err
	}
	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddTapCost().
		AddEffect(abilities.NewCreateTokenEffect(token0_0)).
		Build()
	card.AddAbility(ability0)
	return card, nil
}
