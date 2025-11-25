package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Haunt The Network", NewHauntTheNetwork)
}

// NewHauntTheNetwork creates a Haunt The Network
// {3}{U}{B} - SORCERY
func NewHauntTheNetwork(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Haunt The Network")
	card.ManaCost = "{3}{U}{B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("ThopterColorlessToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainLifeEffect(ArtifactYouControlCount.instance + "where X is the number of artifacts you control")).
		AddEffect(abilities.NewLoseLifeEffect(ArtifactYouControlCount.instance)).
		AddEffect(abilities.NewCreateTokenEffectAmount(token0_0, 2)).
		AddTarget(abilities.NewOpponentTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
