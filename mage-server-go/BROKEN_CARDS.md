# Card Status

**30,600 cards compile & registered** | 93MB binary | Updated: 2025-12-02

## TODO Summary by Category

### Unmapped Abilities (3,424)

| Type                | Count |
| ------------------- | ----- |
| Spell abilities     | 1,876 |
| Activated abilities | 1,548 |

### Missing Triggers (1,442)

| Trigger              | Count |
| -------------------- | ----- |
| LeavesBattlefieldAll | 1,442 |

### Triggered Abilities (1,200+)

| Ability                                          | Count |
| ------------------------------------------------ | ----- |
| EntersBattlefieldTriggeredAbility                | 548   |
| ActivateAsSorceryActivatedAbility                | 136   |
| LoyaltyAbility                                   | 116   |
| AttacksTriggeredAbility                          | 112   |
| BeginningOfUpkeepTriggeredAbility                | 91    |
| ActivateIfConditionActivatedAbility              | 68    |
| DiesSourceTriggeredAbility                       | 66    |
| SpellCastControllerTriggeredAbility              | 46    |
| BeginningOfCombatTriggeredAbility                | 46    |
| BeginningOfEndStepTriggeredAbility               | 44    |
| DealsCombatDamageToAPlayerTriggeredAbility       | 41    |
| SimpleManaAbility                                | 32    |
| EntersBattlefieldOrAttacksSourceTriggeredAbility | 24    |
| LeavesBattlefieldTriggeredAbility                | 19    |
| TurnedFaceUpSourceTriggeredAbility               | 18    |
| DiesCreatureTriggeredAbility                     | 15    |
| CastSourceTriggeredAbility                       | 14    |
| AsEntersBattlefieldAbility                       | 14    |
| EntersBattlefieldThisOrAnotherTriggeredAbility   | 12    |
| SpellCastAllTriggeredAbility                     | 11    |
| LandfallAbility                                  | 10    |
| AttacksAttachedTriggeredAbility                  | 10    |
| Other triggers                                   | ~100  |

### Complex Effect Parameters (3,800+)

| Effect                          | Count |
| ------------------------------- | ----- |
| BoostControlledEffect           | 282   |
| DamageTargetEffect              | 270   |
| GainAbilityAttachedEffect       | 267   |
| ConditionalOneShotEffect        | 244   |
| SearchLibraryPutInHandEffect    | 232   |
| SearchLibraryPutInPlayEffect    | 230   |
| BoostAllEffect                  | 206   |
| BoostSourceEffect               | 192   |
| BoostTargetEffect               | 176   |
| DestroyAllEffect                | 163   |
| GainAbilityControlledEffect     | 137   |
| GainLifeEffect                  | 128   |
| GainAbilityAllEffect            | 103   |
| AddCountersSourceEffect         | 101   |
| GainAbilityTargetEffect         | 75    |
| GainAbilitySourceEffect         | 57    |
| BoostEquippedEffect             | 42    |
| BoostEnchantedEffect            | 41    |
| SearchLibraryPutOnLibraryEffect | 40    |
| AddCountersTargetEffect         | 38    |
| CounterUnlessPaysEffect         | 34    |
| LoseLifeTargetEffect            | 32    |
| ExileAllEffect                  | 26    |
| TapAllEffect                    | 16    |
| LoseLifeSourceControllerEffect  | 14    |
| UntapAllEffect                  | 13    |
| UntapAllControllerEffect        | 12    |
| TapTargetEffect                 | 12    |
| GainControlAllEffect            | 9     |
| GainLifeTargetEffect            | 7     |

### Token Issues (57)

- Token extraction failures

## Commands

```bash
# Verify compilation
go build ./internal/game/cards/generated/...

# Count TODOs
find internal/game/cards/generated -name "*.go" -exec grep -h "TODO:" {} \; | wc -l

# List TODO types
find internal/game/cards/generated -name "*.go" -exec grep -h "TODO:" {} \; | sed 's/.*TODO: //' | sort | uniq -c | sort -rn
```
