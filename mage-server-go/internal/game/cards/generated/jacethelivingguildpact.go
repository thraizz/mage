package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Jace The Living Guildpact", NewJaceTheLivingGuildpact)
}

// NewJaceTheLivingGuildpact creates a Jace The Living Guildpact
// {2}{U}{U} - PLANESWALKER
func NewJaceTheLivingGuildpact(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Jace The Living Guildpact")
	card.ManaCost = "{2}{U}{U}"
	card.Types = []string{"PLANESWALKER"}
	card.Subtypes = []string{"JACE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Loyalty = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		// TODO: ReturnToHandTargetEffect with complex parameters
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryAndPickControllerEffect(2, 1, PutCards.GRAVEYARD, PutCards.TOP_ANY)
	// card.AddAbility(ability1)
	return card, nil
}
