package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Disciple Of Perdition", NewDiscipleOfPerdition)
}

// NewDiscipleOfPerdition creates a Disciple Of Perdition
// {1}{B} - CREATURE
func NewDiscipleOfPerdition(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Disciple Of Perdition")
	card.ManaCost = "{1}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HUMAN", "WARLOCK"}
	card.Power = "1"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
