package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Jeskai Revelation", NewJeskaiRevelation)
}

// NewJeskaiRevelation creates a Jeskai Revelation
// {4}{U}{R}{W} - INSTANT
func NewJeskaiRevelation(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jeskai Revelation")
	card.ManaCost = "{4}{U}{R}{W}"
	card.Types = []string{"INSTANT"}
	card.SetCode = "M21"
	card.Rarity = "common"

	token0_0, err := token.GetToken("MonasteryMentorToken")
	if err != nil {
		return nil, err
	}
	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewReturnToHandTargetEffect()).
		AddEffect(abilities.NewCreateTokenEffect(token0_0, 2)).
		AddEffect(abilities.NewDrawCardsEffect(2)).
		AddEffect(abilities.NewGainLifeEffect(4)).
		AddTarget(abilities.NewAnyTargetFilter()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
