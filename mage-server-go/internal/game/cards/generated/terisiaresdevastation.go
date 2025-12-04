package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Terisiares Devastation", NewTerisiaresDevastation)
}

// NewTerisiaresDevastation creates a Terisiares Devastation
// {X}{2}{B}{B} - SORCERY
func NewTerisiaresDevastation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Terisiares Devastation")
	card.ManaCost = "{X}{2}{B}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("PowerstoneToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: LoseLifeSourceControllerEffect with complex parameters
		// TODO: BoostAllEffect with complex parameters
		AddEffect(abilities.NewCreateTokenEffectAttacking(token0_0, 1, true, false)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
