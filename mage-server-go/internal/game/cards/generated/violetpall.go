package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Violet Pall", NewVioletPall)
}

// NewVioletPall creates a Violet Pall
// {4}{B} - KINDRED INSTANT
func NewVioletPall(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Violet Pall")
	card.ManaCost = "{4}{B}"
	card.Types = []string{"KINDRED", "INSTANT"}
	card.Subtypes = []string{"FAERIE"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("FaerieRogueToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDestroyEffect()).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, 1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
