package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Darth Sidious Sith Lord", NewDarthSidiousSithLord)
}

// NewDarthSidiousSithLord creates a Darth Sidious Sith Lord
// {4}{U}{B}{B}{R} - PLANESWALKER
func NewDarthSidiousSithLord(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Darth Sidious Sith Lord")
	card.ManaCost = "{4}{U}{B}{B}{R}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"SIDIOUS"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: DestroyTargetEffect with complex parameters
		AddEffect(abilities.NewGainControlTargetEffect(abilities.DurationCustom)).
		AddEffect(abilities.NewDamageEffect(7)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
