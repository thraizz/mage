package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Kimahri Valiant Guardian", NewKimahriValiantGuardian)
}

// NewKimahriValiantGuardian creates a Kimahri Valiant Guardian
// {2}{G}{U} - CREATURE
// Vigilance, Vigilance
func NewKimahriValiantGuardian(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Kimahri Valiant Guardian")
	card.ManaCost = "{2}{G}{U}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"CAT", "WARRIOR"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability1)
	return card, nil
}
