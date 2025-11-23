package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
	"github.com/magefree/mage-server-go/internal/game/effects"
)

func init() {
	cards.Register("Cloakwood Hermit", NewCloakwoodHermit)
}

// NewCloakwoodHermit creates a Cloakwood Hermit
// {2}{G} - ENCHANTMENT
func NewCloakwoodHermit(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Cloakwood Hermit")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"ENCHANTMENT"}
	card.Subtypes = []string{"BACKGROUND"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("SquirrelToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, 2, true, false)).
		AddEffect(abilities.NewGrantAbilityEffect(new BeginningOfEndStepTriggeredAbility( TargetController.YOU, new CreateTokenEffect(token0_0, 2, true, false), false, CreaturePutInYourGraveyardCondition.instance ))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}