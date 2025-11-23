package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Defibrillating Current", NewDefibrillatingCurrent)
}

// NewDefibrillatingCurrent creates a Defibrillating Current
// {2/R}{2/W}{2/B} - SORCERY
func NewDefibrillatingCurrent(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Defibrillating Current")
	card.ManaCost = "{2/R}{2/W}{2/B}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(4)).
		AddEffect(abilities.NewGainLifeEffect(2)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
