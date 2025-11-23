package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("The Skullspore Nexus", NewTheSkullsporeNexus)
}

// NewTheSkullsporeNexus creates a The Skullspore Nexus
// {6}{G}{G} - ARTIFACT
func NewTheSkullsporeNexus(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "The Skullspore Nexus")
	card.ManaCost = "{6}{G}{G}"
	card.Types = []string{"ARTIFACT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewCreateTokenEffect(new FungusDinosaurToken(amount))).
		AddEffect(abilities.NewBoostEffect(permanent.getPower().getValue(), 0)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}