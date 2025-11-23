package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Skyshroud Claim", NewSkyshroudClaim)
}

// NewSkyshroudClaim creates a Skyshroud Claim
// {3}{G} - SORCERY
func NewSkyshroudClaim(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Skyshroud Claim")
	card.ManaCost = "{3}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewSearchLibraryPutInPlayEffect(abilities.NewTargetRequirement(0, 1, abilities.NewAnyTargetFilter()), false)).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
