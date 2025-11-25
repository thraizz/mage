package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gideon The Oathsworn", NewGideonTheOathsworn)
}

// NewGideonTheOathsworn creates a Gideon The Oathsworn
// {4}{W}{W} - PLANESWALKER
func NewGideonTheOathsworn(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gideon The Oathsworn")
	card.ManaCost = "{4}{W}{W}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"GIDEON", "SOLDIER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewExileSourceEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
