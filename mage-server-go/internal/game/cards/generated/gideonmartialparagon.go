package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Gideon Martial Paragon", NewGideonMartialParagon)
}

// NewGideonMartialParagon creates a Gideon Martial Paragon
// {4}{W} - PLANESWALKER
// Indestructible
func NewGideonMartialParagon(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Gideon Martial Paragon")
	card.ManaCost = "{4}{W}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"GIDEON", "HUMAN", "SOLDIER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordIndestructible)
	card.AddAbility(ability0)
	ability1, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: UntapAllEffect with complex parameters
		AddEffect(abilities.NewBoostEffect(1, 1)).
		AddEffect(abilities.NewBoostEffect(2, 2)).
		// TODO: TapAllEffect with complex parameters
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability1)
	return card, nil
}
