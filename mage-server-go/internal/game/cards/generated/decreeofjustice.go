package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Decree Of Justice", NewDecreeOfJustice)
}

// NewDecreeOfJustice creates a Decree Of Justice
// {X}{X}{2}{W}{W} - SORCERY
func NewDecreeOfJustice(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Decree Of Justice")
	card.ManaCost = "{X}{X}{2}{W}{W}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	// TODO: Implement triggered ability: CycleTriggeredAbility
	//   - Effect: DecreeOfJusticeCycleEffect()
	// card.AddAbility(ability0)
	token1_0, err := token.GetToken("AngelToken")
	if err != nil {
		return nil, err
	}
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(token1_0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
