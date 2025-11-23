package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("March Of Wretched Sorrow", NewMarchOfWretchedSorrow)
}

// NewMarchOfWretchedSorrow creates a March Of Wretched Sorrow
// {X}{B} - INSTANT
func NewMarchOfWretchedSorrow(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "March Of Wretched Sorrow")
	card.ManaCost = "{X}{B}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDamageEffect(GetXValue.instance)).
		AddEffect(abilities.NewGainLifeEffect(GetXValue.instance)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
