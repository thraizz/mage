package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gideon Jura", NewGideonJura)
}

// NewGideonJura creates a Gideon Jura
// {3}{W}{W} - PLANESWALKER
func NewGideonJura(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gideon Jura")
	card.ManaCost = "{3}{W}{W}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"GIDEON", "HUMAN", "SOLDIER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDestroyEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
