package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Savage Ventmaw", NewSavageVentmaw)
}

// NewSavageVentmaw creates a Savage Ventmaw
// {4}{R}{G} - CREATURE
// Flying
func NewSavageVentmaw(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Savage Ventmaw")
	card.ManaCost = "{4}{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"DRAGON"}
	card.Power = "4"
	card.Toughness = "4"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
