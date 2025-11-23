package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Erinis Gloom Stalker", NewErinisGloomStalker)
}

// NewErinisGloomStalker creates a Erinis Gloom Stalker
// {2}{G} - CREATURE
// Deathtouch
func NewErinisGloomStalker(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Erinis Gloom Stalker")
	card.ManaCost = "{2}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"HALFLING", "RANGER"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordDeathtouch)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	// card.AddAbility(ability1)
	return card, nil
}
