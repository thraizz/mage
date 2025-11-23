package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Tivit Seller Of Secrets", NewTivitSellerOfSecrets)
}

// NewTivitSellerOfSecrets creates a Tivit Seller Of Secrets
// {3}{W}{U}{B} - CREATURE
// Flying
func NewTivitSellerOfSecrets(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Tivit Seller Of Secrets")
	card.ManaCost = "{3}{W}{U}{B}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"SPHINX", "ROGUE"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "6"
	card.Toughness = "6"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordFlying)
	card.AddAbility(ability0)
	return card, nil
}
