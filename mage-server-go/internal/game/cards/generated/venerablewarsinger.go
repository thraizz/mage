package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Venerable Warsinger", NewVenerableWarsinger)
}

// NewVenerableWarsinger creates a Venerable Warsinger
// {1}{R}{W} - CREATURE
// Vigilance, Trample
func NewVenerableWarsinger(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Venerable Warsinger")
	card.ManaCost = "{1}{R}{W}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPIRIT", "CLERIC"}
	card.Power = "3"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	ability1 := abilities.NewKeywordAbility(card.ID, abilities.KeywordTrample)
	card.AddAbility(ability1)
	// TODO: Implement spell ability with unmapped effects
	//   - ReturnFromGraveyardToBattlefieldTargetEffect()
	// card.AddAbility(ability2)
	return card, nil
}
