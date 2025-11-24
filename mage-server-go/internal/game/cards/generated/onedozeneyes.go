package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("One Dozen Eyes", NewOneDozenEyes)
}

// NewOneDozenEyes creates a One Dozen Eyes
// {5}{G} - SORCERY
func NewOneDozenEyes(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "One Dozen Eyes")
	card.ManaCost = "{5}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("InsectToken")
	if err != nil {
		return nil, err
	}
	token0_1, err := token.GetToken("OneDozenEyesBeastToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, 5)).
		AddEffect(abilities.NewCreateTokenEffect(token0_1)).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, 5)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}