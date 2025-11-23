package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Osgir The Reconstructor", NewOsgirTheReconstructor)
}

// NewOsgirTheReconstructor creates a Osgir The Reconstructor
// {2}{R}{W} - CREATURE
// Vigilance
func NewOsgirTheReconstructor(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Osgir The Reconstructor")
	card.ManaCost = "{2}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GIANT", "ARTIFICER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - CreateTokenCopyTargetEffect(player.getId(), null, false, 2)
	// card.AddAbility(ability1)
	ability2 := abilities.NewActivatedAbilityBuilder(card.ID).
		AddEffect(abilities.NewBoostEffect(2, 0)).
		Build()
	card.AddAbility(ability2)
	return card, nil
}
