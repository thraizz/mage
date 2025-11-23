package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Purphoros God Of The Forge", NewPurphorosGodOfTheForge)
}

// NewPurphorosGodOfTheForge creates a Purphoros God Of The Forge
// {3}{R} - ENCHANTMENT CREATURE
// Indestructible
func NewPurphorosGodOfTheForge(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Purphoros God Of The Forge")
	card.ManaCost = "{3}{R}"
	card.Types = []string{"ENCHANTMENT", "CREATURE"}
	card.Subtypes = []string{"GOD"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "5"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordIndestructible)
	card.AddAbility(ability0)
	return card, nil
}
