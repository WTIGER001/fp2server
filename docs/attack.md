
```mermaid
sequenceDiagram
    



```

- Player Clicks on "Attack" icon on a weapon (sword)
- Dialog 
  - Choose character to attack
  - Enter Range (if needed)
- Send Attack Message
- Prompt Defender for Defense Type
  - Select NONE
- Roll Attack
- Resolve Attack
  - Compare attack to DR 15 for Melee
  - Roll Damage (13)
  - Compare damage to Armor (8)
  - Calculate Damage Applied (13-8=5)
  - Degrade Armor SP 13->12
  - Apply to Target Character
  - Update Health
- Show Result To Attacker
- Show Result to Defender
- Mark "Action" as complete