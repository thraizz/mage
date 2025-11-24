package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Stuffy Doll", NewStuffyDoll)
}

// NewStuffyDoll creates a Stuffy Doll
// {5} - ARTIFACT CREATURE
// Indestructible
func NewStuffyDoll(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Stuffy Doll")
	card.ManaCost = "{5}"
	card.Types = []string{"ARTIFACT", "CREATURE"}
	card.Subtypes = []string{"CONSTRUCT"}
	card.Power = "0"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordIndestructible)
	card.AddAbility(ability0)
	// TODO: Implement spell ability with unmapped effects
	//   - ChoosePlayerEffect(Outcome.Damage)
	// card.AddAbility(ability1)
	return card, nil
}