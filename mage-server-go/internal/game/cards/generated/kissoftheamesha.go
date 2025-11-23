package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kiss Of The Amesha", NewKissOfTheAmesha)
}

// NewKissOfTheAmesha creates a Kiss Of The Amesha
// {4}{W}{U} - SORCERY
func NewKissOfTheAmesha(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kiss Of The Amesha")
	card.ManaCost = "{4}{W}{U}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewGainLifeEffect(7)).
		AddEffect(abilities.NewDrawCardsEffect(2)).
		AddTarget(abilities.NewPlayerTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
