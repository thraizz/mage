package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
	"github.com/magefree/mage-server-go/internal/game/token"
)

func init() {
	cards.Register("Rebuild The City", NewRebuildTheCity)
}

// NewRebuildTheCity creates a Rebuild The City
// {3}{B}{R}{G} - SORCERY
// Vigilance
func NewRebuildTheCity(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Rebuild The City")
	card.ManaCost = "{3}{B}{R}{G}"
	card.Types = []string{"SORCERY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - CreateTokenCopyTargetEffect(                 null, CardType.CREATURE, false, 3...)
	//
	// Targets:
	//   - abilities.NewTargetRequirement(1, 1, abilities.NewLandTargetFilter())
	// card.AddAbility(ability1)
	return card, nil
}
