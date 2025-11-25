package generated

import (
	"github.com/google/uuid"
	"github.com/magefree/mage-server-go/internal/game"
	"github.com/magefree/mage-server-go/internal/game/abilities"
	"github.com/magefree/mage-server-go/internal/game/cards"
)

func init() {
	cards.Register("Unscythe Killer Of Kings", NewUnscytheKillerOfKings)
}

// NewUnscytheKillerOfKings creates a Unscythe Killer Of Kings
// {U}{B}{B}{R} - ARTIFACT
func NewUnscytheKillerOfKings(ownerID uuid.UUID, info *cards.CardInfo) (*game.Card, error) {
	card := game.NewCard(ownerID, "Unscythe Killer Of Kings")
	card.ManaCost = "{U}{B}{B}{R}"
	card.Types = []string{"ARTIFACT"}
	card.Subtypes = []string{"EQUIPMENT"}
	card.Supertypes = []string{"LEGENDARY"}
	card.SetCode = "M21"
	card.Rarity = "common"

	ability0, err := abilities.NewSpellAbilityBuilder(card.ID, card.ManaCost).
		AddEffect(abilities.NewBoostEquippedEffect(3, 3)).
		AddEffect(abilities.NewGainAbilityAttachedEffect(abilities.NewKeywordAbility(card.ID, abilities.KeywordFirstStrike), abilities.AttachmentTypeEquipment, abilities.DurationWhileOnBattlefield, "")).
		Build()
	if err != nil {
		return nil, err
	}
	card.AddAbility(ability0)
	return card, nil
}
