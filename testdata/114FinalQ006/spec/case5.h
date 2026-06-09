#ifndef _CASE_H_
#define _CASE_H_
#define SECRET "@SECRET@"

#include "VoyageLog.h"
#include "test.h"

TEST_CASE("Full stress - 4 stops, double hazard, 3 items") {
    {
        auto log = VoyageLog::parse("Iron Depths[3333]18:30?Coral Maze?Jellyfish Bloom@Grid 9 Basin:19:07$Tide Moss:1550173099&Echo Stone:213035152&Grade 3 Ore:1602893089%957201546:21:21?Thermal Spike?Acid Vent@Lava Shelf:22:33$Coconuts:602172146&Pearls:308624087&Kelp Fiber:457466318%1300000000:25:24?Alert 9?Bioluminescent Surge@Silver Lagoon:26:49$Alloy 9:643616802&Crystals:1709503385&Gold Flakes:625470949%1350000000:29:33?Debris Field?Magnetic Field@Thunder Basin:30:42$Echo Stone:1568840169&Gems:1939243944&Gold Flakes:2141599283%1400000000:32:15");

        // ── top-level fields ──
        CHECK(log.shipName == "Iron Depths");
        CHECK(log.equipment == std::vector<int>({3, 3, 3, 3}));
        CHECK((int)log.stops.size() == 4);
        CHECK(log.departMinutes == 1110);
        // ── stop 0: Grid 9 Basin ──
        CHECK(log.stops[0].location == "Grid 9 Basin");
        CHECK(log.stops[0].arriveMinutes == 1147);
        CHECK((int)log.stops[0].hazards.size() == 2);
        CHECK(log.stops[0].hazards[0] == "Coral Maze");
        CHECK(log.stops[0].hazards[1] == "Jellyfish Bloom");
        CHECK((int)log.stops[0].items.size() == 3);
        CHECK(log.stops[0].items[0].name == "Tide Moss");
        CHECK(log.stops[0].items[0].quantity == 1550173099);
        CHECK(log.stops[0].items[1].name == "Echo Stone");
        CHECK(log.stops[0].items[1].quantity == 213035152);
        CHECK(log.stops[0].items[2].name == "Grade 3 Ore");
        CHECK(log.stops[0].items[2].quantity == 1602893089);
        CHECK(log.stops[0].exp == 957201546);
        CHECK(log.stops[0].leaveMinutes == 1281);
        // ── stop 1: Lava Shelf ──
        CHECK(log.stops[1].location == "Lava Shelf");
        CHECK(log.stops[1].arriveMinutes == 1353);
        CHECK((int)log.stops[1].hazards.size() == 2);
        CHECK(log.stops[1].hazards[0] == "Thermal Spike");
        CHECK(log.stops[1].hazards[1] == "Acid Vent");
        CHECK((int)log.stops[1].items.size() == 3);
        CHECK(log.stops[1].items[0].name == "Coconuts");
        CHECK(log.stops[1].items[0].quantity == 602172146);
        CHECK(log.stops[1].items[1].name == "Pearls");
        CHECK(log.stops[1].items[1].quantity == 308624087);
        CHECK(log.stops[1].items[2].name == "Kelp Fiber");
        CHECK(log.stops[1].items[2].quantity == 457466318);
        CHECK(log.stops[1].exp == 1300000000);
        CHECK(log.stops[1].leaveMinutes == 1524);
        // ── stop 2: Silver Lagoon ──
        CHECK(log.stops[2].location == "Silver Lagoon");
        CHECK(log.stops[2].arriveMinutes == 1609);
        CHECK((int)log.stops[2].hazards.size() == 2);
        CHECK(log.stops[2].hazards[0] == "Alert 9");
        CHECK(log.stops[2].hazards[1] == "Bioluminescent Surge");
        CHECK((int)log.stops[2].items.size() == 3);
        CHECK(log.stops[2].items[0].name == "Alloy 9");
        CHECK(log.stops[2].items[0].quantity == 643616802);
        CHECK(log.stops[2].items[1].name == "Crystals");
        CHECK(log.stops[2].items[1].quantity == 1709503385);
        CHECK(log.stops[2].items[2].name == "Gold Flakes");
        CHECK(log.stops[2].items[2].quantity == 625470949);
        CHECK(log.stops[2].exp == 1350000000);
        CHECK(log.stops[2].leaveMinutes == 1773);
        // ── stop 3: Thunder Basin ──
        CHECK(log.stops[3].location == "Thunder Basin");
        CHECK(log.stops[3].arriveMinutes == 1842);
        CHECK((int)log.stops[3].hazards.size() == 2);
        CHECK(log.stops[3].hazards[0] == "Debris Field");
        CHECK(log.stops[3].hazards[1] == "Magnetic Field");
        CHECK((int)log.stops[3].items.size() == 3);
        CHECK(log.stops[3].items[0].name == "Echo Stone");
        CHECK(log.stops[3].items[0].quantity == 1568840169);
        CHECK(log.stops[3].items[1].name == "Gems");
        CHECK(log.stops[3].items[1].quantity == 1939243944);
        CHECK(log.stops[3].items[2].name == "Gold Flakes");
        CHECK(log.stops[3].items[2].quantity == 2141599283);
        CHECK(log.stops[3].exp == 1400000000);
        CHECK(log.stops[3].leaveMinutes == 1935);

        // ── getLocationAt boundaries ──
        CHECK(log.getLocationAt(     0) == "Not yet departed");
        CHECK(log.getLocationAt(  1109) == "Not yet departed");
        CHECK(log.getLocationAt(  1110) == "En route >>>Grid 9 Basin");
        CHECK(log.getLocationAt(  1128) == "En route >>>Grid 9 Basin");
        CHECK(log.getLocationAt(  1146) == "En route >>>Grid 9 Basin");
        CHECK(log.getLocationAt(  1147) == "Grid 9 Basin");
        CHECK(log.getLocationAt(  1214) == "Grid 9 Basin");
        CHECK(log.getLocationAt(  1280) == "Grid 9 Basin");
        CHECK(log.getLocationAt(  1281) == "En route >>>Lava Shelf");
        CHECK(log.getLocationAt(  1317) == "En route >>>Lava Shelf");
        CHECK(log.getLocationAt(  1352) == "En route >>>Lava Shelf");
        CHECK(log.getLocationAt(  1353) == "Lava Shelf");
        CHECK(log.getLocationAt(  1438) == "Lava Shelf");
        CHECK(log.getLocationAt(  1523) == "Lava Shelf");
        CHECK(log.getLocationAt(  1524) == "En route >>>Silver Lagoon");
        CHECK(log.getLocationAt(  1566) == "En route >>>Silver Lagoon");
        CHECK(log.getLocationAt(  1608) == "En route >>>Silver Lagoon");
        CHECK(log.getLocationAt(  1609) == "Silver Lagoon");
        CHECK(log.getLocationAt(  1691) == "Silver Lagoon");
        CHECK(log.getLocationAt(  1772) == "Silver Lagoon");
        CHECK(log.getLocationAt(  1773) == "En route >>>Thunder Basin");
        CHECK(log.getLocationAt(  1807) == "En route >>>Thunder Basin");
        CHECK(log.getLocationAt(  1841) == "En route >>>Thunder Basin");
        CHECK(log.getLocationAt(  1842) == "Thunder Basin");
        CHECK(log.getLocationAt(  1888) == "Thunder Basin");
        CHECK(log.getLocationAt(  1934) == "Thunder Basin");
        CHECK(log.getLocationAt(  1935) == "Returning home");
        CHECK(log.getLocationAt(  2035) == "Returning home");

        // ── getExpAt boundaries ──
        // total = 957201546 + 1300000000 + 1350000000 + 1400000000 = 5007201546 > UINT_MAX
        CHECK(log.getExpAt(1146) == 0);
        CHECK(log.getExpAt(1147) == 957201546LL);
        CHECK(log.getExpAt(1352) == 957201546LL);
        CHECK(log.getExpAt(1353) == 2257201546LL);
        CHECK(log.getExpAt(1608) == 2257201546LL);
        CHECK(log.getExpAt(1609) == 3607201546LL);
        CHECK(log.getExpAt(1841) == 3607201546LL);
        CHECK(log.getExpAt(1842) == 5007201546LL);
        CHECK(log.getExpAt(1935) == 5007201546LL);

        // ── encode() ──
        CHECK(
            log.encode() == "Iron Depths[3333]1110?Coral Maze?Jellyfish Bloom@Grid 9 Basin:1147$Tide Moss:1550173099&Echo Stone:213035152&Grade 3 Ore:1602893089%957201546:1281?Thermal Spike?Acid Vent@Lava Shelf:1353$Coconuts:602172146&Pearls:308624087&Kelp Fiber:457466318%1300000000:1524?Alert 9?Bioluminescent Surge@Silver Lagoon:1609$Alloy 9:643616802&Crystals:1709503385&Gold Flakes:625470949%1350000000:1773?Debris Field?Magnetic Field@Thunder Basin:1842$Echo Stone:1568840169&Gems:1939243944&Gold Flakes:2141599283%1400000000:1935");
    }
    {
        auto log = VoyageLog::parse("Blue 9 Rider[4433]11:45?Kraken Wake?Alert 9@Coral Bay:13:37$Mk5 Fragment:725226997&Tier 8 Fossil:1330973648&Dark Pearl:607371261%686803488:14:09?Jellyfish Bloom?Dark Tide@Shadow Reef:14:47$Frost Shard:1109951537&Tide Moss:961384312&Coconuts:2133288763%1388272068:16:41?Pressure Zone?Bioluminescent Surge@Obsidian Bay:18:32$Magma Rock:1191699301&Tier 8 Fossil:1144607765&Obsidian:301725825%1296773962:20:41?Alert 9?Level 4 Storm@Sector 7 Deep:22:14$Coconuts:1128913215&Class 2 Gem:821508148&Saltstone:11013086%1365886078:24:54");

        // ── top-level fields ──
        CHECK(log.shipName == "Blue 9 Rider");
        CHECK(log.equipment == std::vector<int>({4, 4, 3, 3}));
        CHECK((int)log.stops.size() == 4);
        CHECK(log.departMinutes == 705);
        // ── stop 0: Coral Bay ──
        CHECK(log.stops[0].location == "Coral Bay");
        CHECK(log.stops[0].arriveMinutes == 817);
        CHECK((int)log.stops[0].hazards.size() == 2);
        CHECK(log.stops[0].hazards[0] == "Kraken Wake");
        CHECK(log.stops[0].hazards[1] == "Alert 9");
        CHECK((int)log.stops[0].items.size() == 3);
        CHECK(log.stops[0].items[0].name == "Mk5 Fragment");
        CHECK(log.stops[0].items[0].quantity == 725226997);
        CHECK(log.stops[0].items[1].name == "Tier 8 Fossil");
        CHECK(log.stops[0].items[1].quantity == 1330973648);
        CHECK(log.stops[0].items[2].name == "Dark Pearl");
        CHECK(log.stops[0].items[2].quantity == 607371261);
        CHECK(log.stops[0].exp == 686803488);
        CHECK(log.stops[0].leaveMinutes == 849);
        // ── stop 1: Shadow Reef ──
        CHECK(log.stops[1].location == "Shadow Reef");
        CHECK(log.stops[1].arriveMinutes == 887);
        CHECK((int)log.stops[1].hazards.size() == 2);
        CHECK(log.stops[1].hazards[0] == "Jellyfish Bloom");
        CHECK(log.stops[1].hazards[1] == "Dark Tide");
        CHECK((int)log.stops[1].items.size() == 3);
        CHECK(log.stops[1].items[0].name == "Frost Shard");
        CHECK(log.stops[1].items[0].quantity == 1109951537);
        CHECK(log.stops[1].items[1].name == "Tide Moss");
        CHECK(log.stops[1].items[1].quantity == 961384312);
        CHECK(log.stops[1].items[2].name == "Coconuts");
        CHECK(log.stops[1].items[2].quantity == 2133288763);
        CHECK(log.stops[1].exp == 1388272068);
        CHECK(log.stops[1].leaveMinutes == 1001);
        // ── stop 2: Obsidian Bay ──
        CHECK(log.stops[2].location == "Obsidian Bay");
        CHECK(log.stops[2].arriveMinutes == 1112);
        CHECK((int)log.stops[2].hazards.size() == 2);
        CHECK(log.stops[2].hazards[0] == "Pressure Zone");
        CHECK(log.stops[2].hazards[1] == "Bioluminescent Surge");
        CHECK((int)log.stops[2].items.size() == 3);
        CHECK(log.stops[2].items[0].name == "Magma Rock");
        CHECK(log.stops[2].items[0].quantity == 1191699301);
        CHECK(log.stops[2].items[1].name == "Tier 8 Fossil");
        CHECK(log.stops[2].items[1].quantity == 1144607765);
        CHECK(log.stops[2].items[2].name == "Obsidian");
        CHECK(log.stops[2].items[2].quantity == 301725825);
        CHECK(log.stops[2].exp == 1296773962);
        CHECK(log.stops[2].leaveMinutes == 1241);
        // ── stop 3: Sector 7 Deep ──
        CHECK(log.stops[3].location == "Sector 7 Deep");
        CHECK(log.stops[3].arriveMinutes == 1334);
        CHECK((int)log.stops[3].hazards.size() == 2);
        CHECK(log.stops[3].hazards[0] == "Alert 9");
        CHECK(log.stops[3].hazards[1] == "Level 4 Storm");
        CHECK((int)log.stops[3].items.size() == 3);
        CHECK(log.stops[3].items[0].name == "Coconuts");
        CHECK(log.stops[3].items[0].quantity == 1128913215);
        CHECK(log.stops[3].items[1].name == "Class 2 Gem");
        CHECK(log.stops[3].items[1].quantity == 821508148);
        CHECK(log.stops[3].items[2].name == "Saltstone");
        CHECK(log.stops[3].items[2].quantity == 11013086);
        CHECK(log.stops[3].exp == 1365886078);
        CHECK(log.stops[3].leaveMinutes == 1494);

        // ── getLocationAt boundaries ──
        CHECK(log.getLocationAt(     0) == "Not yet departed");
        CHECK(log.getLocationAt(   704) == "Not yet departed");
        CHECK(log.getLocationAt(   705) == "En route >>>Coral Bay");
        CHECK(log.getLocationAt(   761) == "En route >>>Coral Bay");
        CHECK(log.getLocationAt(   816) == "En route >>>Coral Bay");
        CHECK(log.getLocationAt(   817) == "Coral Bay");
        CHECK(log.getLocationAt(   833) == "Coral Bay");
        CHECK(log.getLocationAt(   848) == "Coral Bay");
        CHECK(log.getLocationAt(   849) == "En route >>>Shadow Reef");
        CHECK(log.getLocationAt(   868) == "En route >>>Shadow Reef");
        CHECK(log.getLocationAt(   886) == "En route >>>Shadow Reef");
        CHECK(log.getLocationAt(   887) == "Shadow Reef");
        CHECK(log.getLocationAt(   944) == "Shadow Reef");
        CHECK(log.getLocationAt(  1000) == "Shadow Reef");
        CHECK(log.getLocationAt(  1001) == "En route >>>Obsidian Bay");
        CHECK(log.getLocationAt(  1056) == "En route >>>Obsidian Bay");
        CHECK(log.getLocationAt(  1111) == "En route >>>Obsidian Bay");
        CHECK(log.getLocationAt(  1112) == "Obsidian Bay");
        CHECK(log.getLocationAt(  1176) == "Obsidian Bay");
        CHECK(log.getLocationAt(  1240) == "Obsidian Bay");
        CHECK(log.getLocationAt(  1241) == "En route >>>Sector 7 Deep");
        CHECK(log.getLocationAt(  1287) == "En route >>>Sector 7 Deep");
        CHECK(log.getLocationAt(  1333) == "En route >>>Sector 7 Deep");
        CHECK(log.getLocationAt(  1334) == "Sector 7 Deep");
        CHECK(log.getLocationAt(  1414) == "Sector 7 Deep");
        CHECK(log.getLocationAt(  1493) == "Sector 7 Deep");
        CHECK(log.getLocationAt(  1494) == "Returning home");
        CHECK(log.getLocationAt(  1594) == "Returning home");

        // ── getExpAt boundaries ──
        // total = 686803488 + 1388272068 + 1296773962 + 1365886078 = 4737735596 > UINT_MAX
        CHECK(log.getExpAt( 816) == 0);
        CHECK(log.getExpAt( 817) == 686803488LL);
        CHECK(log.getExpAt( 886) == 686803488LL);
        CHECK(log.getExpAt( 887) == 2075075556LL);
        CHECK(log.getExpAt(1111) == 2075075556LL);
        CHECK(log.getExpAt(1112) == 3371849518LL);
        CHECK(log.getExpAt(1333) == 3371849518LL);
        CHECK(log.getExpAt(1334) == 4737735596LL);
        CHECK(log.getExpAt(1494) == 4737735596LL);

        // ── encode() ──
        CHECK(
            log.encode() == "Blue 9 Rider[4433]705?Kraken Wake?Alert 9@Coral Bay:817$Mk5 Fragment:725226997&Tier 8 Fossil:1330973648&Dark Pearl:607371261%686803488:849?Jellyfish Bloom?Dark Tide@Shadow Reef:887$Frost Shard:1109951537&Tide Moss:961384312&Coconuts:2133288763%1388272068:1001?Pressure Zone?Bioluminescent Surge@Obsidian Bay:1112$Magma Rock:1191699301&Tier 8 Fossil:1144607765&Obsidian:301725825%1296773962:1241?Alert 9?Level 4 Storm@Sector 7 Deep:1334$Coconuts:1128913215&Class 2 Gem:821508148&Saltstone:11013086%1365886078:1494");
    }
    {
        auto log = VoyageLog::parse("Sea Viper[3313]15:45?Dark Tide?Jellyfish Bloom@Mystery Sea:16:41$Amber:93885964&Class 2 Gem:1963519778&Coconuts:1864575136%777009717:17:42?Bioluminescent Surge?Fog@Platform 11:19:37$Resin 4 Block:1055725&Magma Rock:2003425366&Sea Glass:1759773527%1515362073:22:19?Deep Current?Whirlpool@Lava Shelf:23:52$Fossil Chips:23445418&Batch 0 Kelp:1984216917&Series 6 Amber:1022950740%1241441954:26:41?Jellyfish Bloom?Whirlpool@Blue Beach:27:35$Alloy 9:407598051&Echo Stone:1939387453&Fossil Chips:252849262%1474080636:30:23");

        // ── top-level fields ──
        CHECK(log.shipName == "Sea Viper");
        CHECK(log.equipment == std::vector<int>({3, 3, 1, 3}));
        CHECK((int)log.stops.size() == 4);
        CHECK(log.departMinutes == 945);
        // ── stop 0: Mystery Sea ──
        CHECK(log.stops[0].location == "Mystery Sea");
        CHECK(log.stops[0].arriveMinutes == 1001);
        CHECK((int)log.stops[0].hazards.size() == 2);
        CHECK(log.stops[0].hazards[0] == "Dark Tide");
        CHECK(log.stops[0].hazards[1] == "Jellyfish Bloom");
        CHECK((int)log.stops[0].items.size() == 3);
        CHECK(log.stops[0].items[0].name == "Amber");
        CHECK(log.stops[0].items[0].quantity == 93885964);
        CHECK(log.stops[0].items[1].name == "Class 2 Gem");
        CHECK(log.stops[0].items[1].quantity == 1963519778);
        CHECK(log.stops[0].items[2].name == "Coconuts");
        CHECK(log.stops[0].items[2].quantity == 1864575136);
        CHECK(log.stops[0].exp == 777009717);
        CHECK(log.stops[0].leaveMinutes == 1062);
        // ── stop 1: Platform 11 ──
        CHECK(log.stops[1].location == "Platform 11");
        CHECK(log.stops[1].arriveMinutes == 1177);
        CHECK((int)log.stops[1].hazards.size() == 2);
        CHECK(log.stops[1].hazards[0] == "Bioluminescent Surge");
        CHECK(log.stops[1].hazards[1] == "Fog");
        CHECK((int)log.stops[1].items.size() == 3);
        CHECK(log.stops[1].items[0].name == "Resin 4 Block");
        CHECK(log.stops[1].items[0].quantity == 1055725);
        CHECK(log.stops[1].items[1].name == "Magma Rock");
        CHECK(log.stops[1].items[1].quantity == 2003425366);
        CHECK(log.stops[1].items[2].name == "Sea Glass");
        CHECK(log.stops[1].items[2].quantity == 1759773527);
        CHECK(log.stops[1].exp == 1515362073);
        CHECK(log.stops[1].leaveMinutes == 1339);
        // ── stop 2: Lava Shelf ──
        CHECK(log.stops[2].location == "Lava Shelf");
        CHECK(log.stops[2].arriveMinutes == 1432);
        CHECK((int)log.stops[2].hazards.size() == 2);
        CHECK(log.stops[2].hazards[0] == "Deep Current");
        CHECK(log.stops[2].hazards[1] == "Whirlpool");
        CHECK((int)log.stops[2].items.size() == 3);
        CHECK(log.stops[2].items[0].name == "Fossil Chips");
        CHECK(log.stops[2].items[0].quantity == 23445418);
        CHECK(log.stops[2].items[1].name == "Batch 0 Kelp");
        CHECK(log.stops[2].items[1].quantity == 1984216917);
        CHECK(log.stops[2].items[2].name == "Series 6 Amber");
        CHECK(log.stops[2].items[2].quantity == 1022950740);
        CHECK(log.stops[2].exp == 1241441954);
        CHECK(log.stops[2].leaveMinutes == 1601);
        // ── stop 3: Blue Beach ──
        CHECK(log.stops[3].location == "Blue Beach");
        CHECK(log.stops[3].arriveMinutes == 1655);
        CHECK((int)log.stops[3].hazards.size() == 2);
        CHECK(log.stops[3].hazards[0] == "Jellyfish Bloom");
        CHECK(log.stops[3].hazards[1] == "Whirlpool");
        CHECK((int)log.stops[3].items.size() == 3);
        CHECK(log.stops[3].items[0].name == "Alloy 9");
        CHECK(log.stops[3].items[0].quantity == 407598051);
        CHECK(log.stops[3].items[1].name == "Echo Stone");
        CHECK(log.stops[3].items[1].quantity == 1939387453);
        CHECK(log.stops[3].items[2].name == "Fossil Chips");
        CHECK(log.stops[3].items[2].quantity == 252849262);
        CHECK(log.stops[3].exp == 1474080636);
        CHECK(log.stops[3].leaveMinutes == 1823);

        // ── getLocationAt boundaries ──
        CHECK(log.getLocationAt(     0) == "Not yet departed");
        CHECK(log.getLocationAt(   944) == "Not yet departed");
        CHECK(log.getLocationAt(   945) == "En route >>>Mystery Sea");
        CHECK(log.getLocationAt(   973) == "En route >>>Mystery Sea");
        CHECK(log.getLocationAt(  1000) == "En route >>>Mystery Sea");
        CHECK(log.getLocationAt(  1001) == "Mystery Sea");
        CHECK(log.getLocationAt(  1031) == "Mystery Sea");
        CHECK(log.getLocationAt(  1061) == "Mystery Sea");
        CHECK(log.getLocationAt(  1062) == "En route >>>Platform 11");
        CHECK(log.getLocationAt(  1119) == "En route >>>Platform 11");
        CHECK(log.getLocationAt(  1176) == "En route >>>Platform 11");
        CHECK(log.getLocationAt(  1177) == "Platform 11");
        CHECK(log.getLocationAt(  1258) == "Platform 11");
        CHECK(log.getLocationAt(  1338) == "Platform 11");
        CHECK(log.getLocationAt(  1339) == "En route >>>Lava Shelf");
        CHECK(log.getLocationAt(  1385) == "En route >>>Lava Shelf");
        CHECK(log.getLocationAt(  1431) == "En route >>>Lava Shelf");
        CHECK(log.getLocationAt(  1432) == "Lava Shelf");
        CHECK(log.getLocationAt(  1516) == "Lava Shelf");
        CHECK(log.getLocationAt(  1600) == "Lava Shelf");
        CHECK(log.getLocationAt(  1601) == "En route >>>Blue Beach");
        CHECK(log.getLocationAt(  1628) == "En route >>>Blue Beach");
        CHECK(log.getLocationAt(  1654) == "En route >>>Blue Beach");
        CHECK(log.getLocationAt(  1655) == "Blue Beach");
        CHECK(log.getLocationAt(  1739) == "Blue Beach");
        CHECK(log.getLocationAt(  1822) == "Blue Beach");
        CHECK(log.getLocationAt(  1823) == "Returning home");
        CHECK(log.getLocationAt(  1923) == "Returning home");

        // ── getExpAt boundaries ──
        // total = 777009717 + 1515362073 + 1241441954 + 1474080636 = 5007894380 > UINT_MAX
        CHECK(log.getExpAt(1000) == 0);
        CHECK(log.getExpAt(1001) == 777009717LL);
        CHECK(log.getExpAt(1176) == 777009717LL);
        CHECK(log.getExpAt(1177) == 2292371790LL);
        CHECK(log.getExpAt(1431) == 2292371790LL);
        CHECK(log.getExpAt(1432) == 3533813744LL);
        CHECK(log.getExpAt(1654) == 3533813744LL);
        CHECK(log.getExpAt(1655) == 5007894380LL);
        CHECK(log.getExpAt(1823) == 5007894380LL);

        // ── encode() ──
        CHECK(
            log.encode() == "Sea Viper[3313]945?Dark Tide?Jellyfish Bloom@Mystery Sea:1001$Amber:93885964&Class 2 Gem:1963519778&Coconuts:1864575136%777009717:1062?Bioluminescent Surge?Fog@Platform 11:1177$Resin 4 Block:1055725&Magma Rock:2003425366&Sea Glass:1759773527%1515362073:1339?Deep Current?Whirlpool@Lava Shelf:1432$Fossil Chips:23445418&Batch 0 Kelp:1984216917&Series 6 Amber:1022950740%1241441954:1601?Jellyfish Bloom?Whirlpool@Blue Beach:1655$Alloy 9:407598051&Echo Stone:1939387453&Fossil Chips:252849262%1474080636:1823");
    }
}

#endif // _CASE_H_
