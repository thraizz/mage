package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Aven Fateshaper", NewAvenFateshaper)
}

// NewAvenFateshaper creates a Aven Fateshaper
// {6}{U} - CREATURE
// Flying
func NewAvenFateshaper(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Aven Fateshaper")
	card.ManaCost = "{6}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"BIRD", "WIZARD"}
	card.Power = "4"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - LookLibraryControllerEffect(4)
	// card.AddAbility(ability1)
	return card, nil
}
