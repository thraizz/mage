package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Celestial Colonnade", NewCelestialColonnade)
}

// NewCelestialColonnade creates a Celestial Colonnade
//   - LAND
//
// Flying, Vigilance
func NewCelestialColonnade(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Celestial Colonnade")
	card.ManaCost = ""
	card.Types = []string{"LAND"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	ability2 := abilities.BuildSimpleManaAbility(card.ID, "U")
	card.AddAbility(ability2)
	ability3 := abilities.BuildSimpleManaAbility(card.ID, "W")
	card.AddAbility(ability3)
	return card, nil
}
