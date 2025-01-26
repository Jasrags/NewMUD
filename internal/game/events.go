package game

import ee "github.com/vansante/go-event-emitter"

const (
	TestEvent ee.EventType = "test"

	EventRoomChannelReceive ee.EventType = "room_channel_receive"
	EventRoomMobEnter       ee.EventType = "room_mob_enter"
	EventRoomMobLeave       ee.EventType = "room_mob_leave"
	EventRoomCharacterEnter ee.EventType = "room_character_enter"
	EventRoomCharacterLeave ee.EventType = "room_character_leave"
	// EventRoomReady          ee.EventType = "room_ready"
	// EventRoomSpawn          ee.EventType = "room_spawn"
	// EventRoomUpdate         ee.EventType = "room_update"

	EventPlayerEnterRoom ee.EventType = "player_enter_room"
)

type (
	RoomCharacterEnter struct {
		Character *Character
		Room      *Room
		PrevRoom  *Room
	}
	RoomCharacterLeave struct {
		Character *Character
		Room      *Room
		NextRoom  *Room
	}
	RoomMobEnter struct {
		Mob      *Mob
		Room     *Room
		PrevRoom *Room
	}
	RoomMobLeave struct {
		Mob      *Mob
		Room     *Room
		NextRoom *Room
	}
	PlayerEnterRoom struct {
		Character *Character
		Room      *Room
	}
)

// /Users/jrags/Code/Ranvier/core/src/Area.js
//   69,5:    * @fires Area#roomAdded
//   87,5:    * @fires Area#roomRemoved
//   151,5:    * @fires Room#updateTick
//   152,5:    * @fires Npc#updateTick

// /Users/jrags/Code/Ranvier/core/src/AreaManager.js
//   51,5:    * @fires Area#updateTick

// /Users/jrags/Code/Ranvier/core/src/Channel.js
//   50,5:    * @fires GameEntity#channelReceive

// /Users/jrags/Code/Ranvier/core/src/Character.js
//   194,5:    * @fires Character#attributeUpdate
//   234,5:    * @fires Character#combatStart
//   274,5:    * @fires Character#combatantAdded
//   292,5:    * @fires Character#combatantRemoved
//   293,5:    * @fires Character#combatEnd
//   357,5:    * @fires Character#equip
//   358,5:    * @fires Item#equip
//   394,5:    * @fires Item#unequip
//   395,5:    * @fires Character#unequip
//   504,5:    * @fires Character#unfollowed
//   518,5:    * @fires Character#gainedFollower
//   532,5:    * @fires Character#lostFollower

// /Users/jrags/Code/Ranvier/core/src/Damage.js
//   52,5:    * @fires Character#hit
//   53,5:    * @fires Character#damaged

// /Users/jrags/Code/Ranvier/core/src/Effect.js
//   137,5:    * @fires Effect#effectActivated
//   158,5:    * @fires Effect#effectDeactivated
//   174,5:    * @fires Effect#remove

// /Users/jrags/Code/Ranvier/core/src/EffectList.js
//   87,5:    * @fires Effect#effectAdded
//   88,5:    * @fires Effect#effectStackAdded
//   89,5:    * @fires Effect#effectRefreshed
//   90,5:    * @fires Character#effectAdded
//   144,5:    * @fires Character#effectRemoved

// /Users/jrags/Code/Ranvier/core/src/GameServer.js
//   9,5:    * @fires GameServer#startup
//   20,5:    * @fires GameServer#shutdown

// /Users/jrags/Code/Ranvier/core/src/Heal.js
//   13,5:    * @fires Character#heal
//   14,5:    * @fires Character#healed

// /Users/jrags/Code/Ranvier/core/src/ItemManager.js
//   36,5:    * @fires Item#updateTick

// /Users/jrags/Code/Ranvier/core/src/Metadatable.js
//   24,5:    * @fires Metadatable#metadataUpdate

// /Users/jrags/Code/Ranvier/core/src/Npc.js
//   48,5:    * @fires Room#npcLeave
//   49,5:    * @fires Room#npcEnter
//   50,5:    * @fires Npc#enterRoom

// /Users/jrags/Code/Ranvier/core/src/Player.js
//   131,5:    * @fires Room#playerLeave
//   132,5:    * @fires Room#playerEnter
//   133,5:    * @fires Player#enterRoom

// /Users/jrags/Code/Ranvier/core/src/PlayerManager.js
//   141,5:    * @fires Player#save
//   157,5:    * @fires Player#saved
//   166,5:    * @fires Player#updateTick

// /Users/jrags/Code/Ranvier/core/src/Quest.js
//   59,5:    * @fires Quest#turn-in-ready
//   60,5:    * @fires Quest#progress
//   125,5:    * @fires Quest#complete

// /Users/jrags/Code/Ranvier/core/src/Room.js
//   325,5:    * @fires Npc#spawn

// /Users/jrags/Code/Ranvier/core/src/Area.js
//   79,7:      * @event Area#roomAdded
//   93,7:      * @event Area#roomRemoved
//   158,9:        * @event Room#updateTick
//   166,9:        * @event Npc#updateTick
//   182,9:        * @event Room#ready

// /Users/jrags/Code/Ranvier/core/src/AreaManager.js
//   57,9:        * @event Area#updateTick

// /Users/jrags/Code/Ranvier/core/src/Channel.js
//   96,9:        * @event GameEntity#channelReceive

// /Users/jrags/Code/Ranvier/core/src/Character.js
//   132,5:    * @event Character#attributeUpdate
//   242,9:        * @event Character#combatStart
//   284,7:      * @event Character#combatantAdded
//   304,7:      * @event Character#combatantRemoved
//   311,9:        * @event Character#combatEnd
//   377,7:      * @event Item#equip
//   382,7:      * @event Character#equip
//   407,7:      * @event Item#unequip
//   412,7:      * @event Character#unequip
//   496,7:      * @event Character#followed
//   509,7:      * @event Character#unfollowed
//   524,7:      * @event Character#gainedFollower
//   538,7:      * @event Character#lostFollower

// /Users/jrags/Code/Ranvier/core/src/Damage.js
//   61,9:        * @event Character#hit
//   69,9:        * @event Character#damaged

// /Users/jrags/Code/Ranvier/core/src/Effect.js
//   150,7:      * @event Effect#effectActivated
//   166,7:      * @event Effect#effectDeactivated
//   178,7:      * @event Effect#remove

// /Users/jrags/Code/Ranvier/core/src/EffectList.js
//   103,13:            * @event Effect#effectStackAdded
//   112,13:            * @event Effect#effectRefreshed
//   129,7:      * @event Effect#effectAdded
//   133,7:      * @event Character#effectAdded
//   154,7:      * @event Character#effectRemoved

// /Users/jrags/Code/Ranvier/core/src/GameServer.js
//   13,7:      * @event GameServer#startup
//   24,7:      * @event GameServer#shutdown

// /Users/jrags/Code/Ranvier/core/src/Heal.js
//   22,9:        * @event Character#heal
//   30,7:      * @event Character#healed

// /Users/jrags/Code/Ranvier/core/src/ItemManager.js
//   41,9:        * @event Item#updateTick

// /Users/jrags/Code/Ranvier/core/src/Metadatable.js
//   47,6:     * @event Metadatable#metadataUpdate

// /Users/jrags/Code/Ranvier/core/src/Npc.js
//   56,9:        * @event Room#npcLeave
//   70,7:      * @event Room#npcEnter
//   76,7:      * @event Npc#enterRoom

// /Users/jrags/Code/Ranvier/core/src/Player.js
//   139,9:        * @event Room#playerLeave
//   153,7:      * @event Room#playerEnter
//   159,7:      * @event Player#enterRoom

// /Users/jrags/Code/Ranvier/core/src/PlayerManager.js
//   151,7:      * @event Player#saved
//   171,9:        * @event Player#updateTick

// /Users/jrags/Code/Ranvier/core/src/Quest.js
//   70,11:          * @event Quest#turn-in-ready
//   78,7:      * @event Quest#progress
//   129,7:      * @event Quest#complete

// /Users/jrags/Code/Ranvier/core/src/Room.js
//   316,7:      * @event Item#spawn
//   337,7:      * @event Npc#spawn
//   349,7:      * @event Room#spawn
