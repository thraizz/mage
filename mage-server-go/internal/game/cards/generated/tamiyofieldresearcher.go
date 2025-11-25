package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tamiyo Field Researcher", NewTamiyoFieldResearcher)
}

// NewTamiyoFieldResearcher creates a Tamiyo Field Researcher
// {1}{G}{W}{U} - PLANESWALKER
func NewTamiyoFieldResearcher(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tamiyo Field Researcher")
	card.ManaCost = "{1}{G}{W}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"TAMIYO"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewDrawCardsEffect(1)).
		AddEffect(abilities.NewTapEffect()).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
