package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Summon The School", NewSummonTheSchool)
}

// NewSummonTheSchool creates a Summon The School
// {3}{W} - KINDRED SORCERY
func NewSummonTheSchool(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Summon The School")
	card.ManaCost = "{3}{W}"
	card.Types = []string{"KINDRED", "SORCERY"}
	card.Subtypes = []string{"MERFOLK"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("MerfolkWizardToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, 2)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}