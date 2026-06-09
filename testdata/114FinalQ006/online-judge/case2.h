#ifndef _CASE_H_
#define _CASE_H_
#define SECRET "5c16ad0f149ab0c559acacbf013ebfa4b928efbe1f5acd14a97f59945cc5d7f9"

#include "VoyageLog.h"
#include "test.h"

TEST_CASE("Parse - hazards and multi-item") {
    {
        auto log = VoyageLog::parse("Abyss Rider[4442]2:45?Fog@Ember Isle:4:11$Frost Shard:3303&Sea Glass:3891%64785:6:36?Ice Drift@Twin Peaks:8:26$Gold Flakes:7801&Sea Glass:2104%55800:11:09");

        CHECK(log.shipName == "Abyss Rider");
        CHECK(log.equipment == std::vector<int>({4, 4, 4, 2}));
        CHECK((int)log.stops.size() == 2);
        CHECK(log.departMinutes == 165);

        // ── Stop 0: Ember Isle ──
        CHECK((int)log.stops[0].hazards.size() == 1);
        CHECK(log.stops[0].hazards[0] == "Fog");
        CHECK(log.stops[0].location == "Ember Isle");
        CHECK(log.stops[0].arriveMinutes == 251);
        CHECK((int)log.stops[0].items.size() == 2);
        CHECK(log.stops[0].items[0].name == "Frost Shard");
        CHECK(log.stops[0].items[0].quantity == 3303);
        CHECK(log.stops[0].items[1].name == "Sea Glass");
        CHECK(log.stops[0].items[1].quantity == 3891);
        CHECK(log.stops[0].exp == 64785);
        CHECK(log.stops[0].leaveMinutes == 396);

        // ── Stop 1: Twin Peaks ──
        CHECK((int)log.stops[1].hazards.size() == 1);
        CHECK(log.stops[1].hazards[0] == "Ice Drift");
        CHECK(log.stops[1].location == "Twin Peaks");
        CHECK(log.stops[1].arriveMinutes == 506);
        CHECK((int)log.stops[1].items.size() == 2);
        CHECK(log.stops[1].items[0].name == "Gold Flakes");
        CHECK(log.stops[1].items[0].quantity == 7801);
        CHECK(log.stops[1].items[1].name == "Sea Glass");
        CHECK(log.stops[1].items[1].quantity == 2104);
        CHECK(log.stops[1].exp == 55800);
        CHECK(log.stops[1].leaveMinutes == 669);
    }
    {
        auto log = VoyageLog::parse("Storm Chaser[1234]9?Fog@Lava Shelf:9:52$Gold Flakes:191&Saltstone:2897%27918:10:47?Debris Field@Obsidian Bay:11:42$Ores:3245&Obsidian:7139%31168:12:43");

        CHECK(log.shipName == "Storm Chaser");
        CHECK(log.equipment == std::vector<int>({1, 2, 3, 4}));
        CHECK((int)log.stops.size() == 2);
        CHECK(log.departMinutes == 540);

        // ── Stop 0: Lava Shelf ──
        CHECK((int)log.stops[0].hazards.size() == 1);
        CHECK(log.stops[0].hazards[0] == "Fog");
        CHECK(log.stops[0].location == "Lava Shelf");
        CHECK(log.stops[0].arriveMinutes == 592);
        CHECK((int)log.stops[0].items.size() == 2);
        CHECK(log.stops[0].items[0].name == "Gold Flakes");
        CHECK(log.stops[0].items[0].quantity == 191);
        CHECK(log.stops[0].items[1].name == "Saltstone");
        CHECK(log.stops[0].items[1].quantity == 2897);
        CHECK(log.stops[0].exp == 27918);
        CHECK(log.stops[0].leaveMinutes == 647);

        // ── Stop 1: Obsidian Bay ──
        CHECK((int)log.stops[1].hazards.size() == 1);
        CHECK(log.stops[1].hazards[0] == "Debris Field");
        CHECK(log.stops[1].location == "Obsidian Bay");
        CHECK(log.stops[1].arriveMinutes == 702);
        CHECK((int)log.stops[1].items.size() == 2);
        CHECK(log.stops[1].items[0].name == "Ores");
        CHECK(log.stops[1].items[0].quantity == 3245);
        CHECK(log.stops[1].items[1].name == "Obsidian");
        CHECK(log.stops[1].items[1].quantity == 7139);
        CHECK(log.stops[1].exp == 31168);
        CHECK(log.stops[1].leaveMinutes == 763);
    }
    {
        auto log = VoyageLog::parse("Iron Depths[3233]8:30?Jellyfish Bloom@Sunken Arch:9:07$Gold Flakes:9386&Frost Shard:5672%45804:10:26?Acid Vent@Jade Grotto:11:48$Echo Stone:5655&Iron Shards:1820%25972:12:58");

        CHECK(log.shipName == "Iron Depths");
        CHECK(log.equipment == std::vector<int>({3, 2, 3, 3}));
        CHECK((int)log.stops.size() == 2);
        CHECK(log.departMinutes == 510);

        // ── Stop 0: Sunken Arch ──
        CHECK((int)log.stops[0].hazards.size() == 1);
        CHECK(log.stops[0].hazards[0] == "Jellyfish Bloom");
        CHECK(log.stops[0].location == "Sunken Arch");
        CHECK(log.stops[0].arriveMinutes == 547);
        CHECK((int)log.stops[0].items.size() == 2);
        CHECK(log.stops[0].items[0].name == "Gold Flakes");
        CHECK(log.stops[0].items[0].quantity == 9386);
        CHECK(log.stops[0].items[1].name == "Frost Shard");
        CHECK(log.stops[0].items[1].quantity == 5672);
        CHECK(log.stops[0].exp == 45804);
        CHECK(log.stops[0].leaveMinutes == 626);

        // ── Stop 1: Jade Grotto ──
        CHECK((int)log.stops[1].hazards.size() == 1);
        CHECK(log.stops[1].hazards[0] == "Acid Vent");
        CHECK(log.stops[1].location == "Jade Grotto");
        CHECK(log.stops[1].arriveMinutes == 708);
        CHECK((int)log.stops[1].items.size() == 2);
        CHECK(log.stops[1].items[0].name == "Echo Stone");
        CHECK(log.stops[1].items[0].quantity == 5655);
        CHECK(log.stops[1].items[1].name == "Iron Shards");
        CHECK(log.stops[1].items[1].quantity == 1820);
        CHECK(log.stops[1].exp == 25972);
        CHECK(log.stops[1].leaveMinutes == 778);
    }
}

#endif // _CASE_H_
