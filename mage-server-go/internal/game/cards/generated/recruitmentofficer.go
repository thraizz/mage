package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Recruitment Officer", NewRecruitmentOfficer)
}

// NewRecruitmentOfficer creates a Recruitment Officer
// {W} - CREATURE
func NewRecruitmentOfficer(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Recruitment Officer")
	card.ManaCost = "{W}"
	card.Types = []string{"CREATURE"}
	card.Power = "2"
	card.Toughness = "1"
	card.SetCode = "M21"
	card.Rarity = "common"

	return card, nil
}
