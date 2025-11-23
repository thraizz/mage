package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("The Battle Of Geonosis", NewTheBattleOfGeonosis)
}

// NewTheBattleOfGeonosis creates a The Battle Of Geonosis
// {X}{X}{R}{R} - SORCERY
func NewTheBattleOfGeonosis(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Battle Of Geonosis")
	card.ManaCost = "{X}{X}{R}{R}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(new IntPlusDynamicValue(1, GetXValue.instance))).
		AddEffect(abilities.NewDamageEffect(new IntPlusDynamicValue(1, GetXValue.instance))).
		AddEffect(abilities.NewBoostEffect(GetXValue.instance, StaticValue.get(0))).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}