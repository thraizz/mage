package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Tinkers Tote", NewTinkersTote)
}

// NewTinkersTote creates a Tinkers Tote
// {2}{W} - ARTIFACT
func NewTinkersTote(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tinkers Tote")
	card.ManaCost = "{2}{W}"
	card.Types = []string{"ARTIFACT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddSacrificeSourceCost().
		AddEffect(abilities.NewGainLifeEffect(3)).
		Build()
	card.AddAbility(ability0)
	token1_0, err := token.GetToken("GnomeToken")
	if err != nil {
		return nil, err
	}
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token1_0, 2)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
