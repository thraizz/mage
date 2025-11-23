package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Redshift Rocketeer Chief", NewRedshiftRocketeerChief)
}

// NewRedshiftRocketeerChief creates a Redshift Rocketeer Chief
// {R}{G} - CREATURE
// Vigilance
func NewRedshiftRocketeerChief(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Redshift Rocketeer Chief")
	card.ManaCost = "{R}{G}"
	card.Types = []string{"CREATURE"}
	card.Subtypes = []string{"GOBLIN", "PILOT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.Power = "2"
	card.Toughness = "3"
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0 := abilities.NewKeywordAbility(card.ID, abilities.KeywordVigilance)
	card.AddAbility(ability0)
	return card, nil
}
