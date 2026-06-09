#ifndef _CASE_H_
#define _CASE_H_
#define SECRET "@SECRET@"

#include "VoyageLog.h"
#include "test.h"

TEST_CASE("getExpAt + encode, large values") {
    {
        auto log = VoyageLog::parse(
            "Leviathan[4111]23?Deep Current@Sunken Arch:23:42$Batch 0 Kelp:1434009068&Grade 3 Ore:1811080350%515381578:24:17?Thermal Vent@Crystal Cove:24:58$Pearls:1416305517&Amber:1257600446%711414429:27:44?Pressure Wave@Twin Peaks:28:37$Ores:1086811674&Tier 8 Fossil:1569175984%142203829:29:55");

        // ── top-level fields ──
        CHECK(log.shipName == "Leviathan");
        CHECK(log.equipment == std::vector<int>({4, 1, 1, 1}));
        CHECK((int)log.stops.size() == 3);
        CHECK(log.departMinutes == 1380);
        // ── stop 0: Sunken Arch ──
        CHECK((int)log.stops[0].hazards.size() == 1);
        CHECK(log.stops[0].hazards[0] == "Deep Current");
        CHECK(log.stops[0].location == "Sunken Arch");
        CHECK(log.stops[0].arriveMinutes == 1422);
        CHECK((int)log.stops[0].items.size() == 2);
        CHECK(log.stops[0].items[0].name == "Batch 0 Kelp");
        CHECK(log.stops[0].items[0].quantity == 1434009068);
        CHECK(log.stops[0].items[1].name == "Grade 3 Ore");
        CHECK(log.stops[0].items[1].quantity == 1811080350);
        CHECK(log.stops[0].exp == 515381578);
        CHECK(log.stops[0].leaveMinutes == 1457);
        // ── stop 1: Crystal Cove ──
        CHECK((int)log.stops[1].hazards.size() == 1);
        CHECK(log.stops[1].hazards[0] == "Thermal Vent");
        CHECK(log.stops[1].location == "Crystal Cove");
        CHECK(log.stops[1].arriveMinutes == 1498);
        CHECK((int)log.stops[1].items.size() == 2);
        CHECK(log.stops[1].items[0].name == "Pearls");
        CHECK(log.stops[1].items[0].quantity == 1416305517);
        CHECK(log.stops[1].items[1].name == "Amber");
        CHECK(log.stops[1].items[1].quantity == 1257600446);
        CHECK(log.stops[1].exp == 711414429);
        CHECK(log.stops[1].leaveMinutes == 1664);
        // ── stop 2: Twin Peaks ──
        CHECK((int)log.stops[2].hazards.size() == 1);
        CHECK(log.stops[2].hazards[0] == "Pressure Wave");
        CHECK(log.stops[2].location == "Twin Peaks");
        CHECK(log.stops[2].arriveMinutes == 1717);
        CHECK((int)log.stops[2].items.size() == 2);
        CHECK(log.stops[2].items[0].name == "Ores");
        CHECK(log.stops[2].items[0].quantity == 1086811674);
        CHECK(log.stops[2].items[1].name == "Tier 8 Fossil");
        CHECK(log.stops[2].items[1].quantity == 1569175984);
        CHECK(log.stops[2].exp == 142203829);
        CHECK(log.stops[2].leaveMinutes == 1795);

        // ── getLocationAt boundaries ──
        CHECK(log.getLocationAt(     0) == "Not yet departed");
        CHECK(log.getLocationAt(  1379) == "Not yet departed");
        CHECK(log.getLocationAt(  1380) == "En route >>>Sunken Arch");
        CHECK(log.getLocationAt(  1401) == "En route >>>Sunken Arch");
        CHECK(log.getLocationAt(  1421) == "En route >>>Sunken Arch");
        CHECK(log.getLocationAt(  1422) == "Sunken Arch");
        CHECK(log.getLocationAt(  1439) == "Sunken Arch");
        CHECK(log.getLocationAt(  1456) == "Sunken Arch");
        CHECK(log.getLocationAt(  1457) == "En route >>>Crystal Cove");
        CHECK(log.getLocationAt(  1477) == "En route >>>Crystal Cove");
        CHECK(log.getLocationAt(  1497) == "En route >>>Crystal Cove");
        CHECK(log.getLocationAt(  1498) == "Crystal Cove");
        CHECK(log.getLocationAt(  1581) == "Crystal Cove");
        CHECK(log.getLocationAt(  1663) == "Crystal Cove");
        CHECK(log.getLocationAt(  1664) == "En route >>>Twin Peaks");
        CHECK(log.getLocationAt(  1690) == "En route >>>Twin Peaks");
        CHECK(log.getLocationAt(  1716) == "En route >>>Twin Peaks");
        CHECK(log.getLocationAt(  1717) == "Twin Peaks");
        CHECK(log.getLocationAt(  1756) == "Twin Peaks");
        CHECK(log.getLocationAt(  1794) == "Twin Peaks");
        CHECK(log.getLocationAt(  1795) == "Returning home");
        CHECK(log.getLocationAt(  1895) == "Returning home");

        // ── getExpAt boundaries ──
        CHECK(log.getExpAt(1421) == 0u);
        CHECK(log.getExpAt(1421) == 0);
        CHECK(log.getExpAt(1422) == 515381578);
        CHECK(log.getExpAt(1497) == 515381578);
        CHECK(log.getExpAt(1498) == 1226796007);
        CHECK(log.getExpAt(1716) == 1226796007);
        CHECK(log.getExpAt(1717) == 1368999836);
        CHECK(log.getExpAt(1795) == 1368999836);

        // ── encode() outputs pure minutes, not HH:MM ──
        CHECK(
            log.encode() ==
            "Leviathan[4111]1380?Deep Current@Sunken Arch:1422$Batch 0 Kelp:1434009068&Grade 3 Ore:1811080350%515381578:1457?Thermal Vent@Crystal Cove:1498$Pearls:1416305517&Amber:1257600446%711414429:1664?Pressure Wave@Twin Peaks:1717$Ores:1086811674&Tier 8 Fossil:1569175984%142203829:1795");
    }
    {
        auto log =
            VoyageLog::parse("Nova Shell[4313]23?Fog Bank@Sunken Arch:24:39$Gems:584633497&Crystals:202209422%456423512:26:02?Whirlpool@Hex 0 Atoll:27:28$Dark Pearl:1175006141&Amber:1848890585%415252392:29:47?Dark Current@Site 8 Hollow:30:31$Resin 4 Block:1781911614&Coral:1508751711%317796632:33:10");

        // ── top-level fields ──
        CHECK(log.shipName == "Nova Shell");
        CHECK(log.equipment == std::vector<int>({4, 3, 1, 3}));
        CHECK((int)log.stops.size() == 3);
        CHECK(log.departMinutes == 1380);
        // ── stop 0: Sunken Arch ──
        CHECK((int)log.stops[0].hazards.size() == 1);
        CHECK(log.stops[0].hazards[0] == "Fog Bank");
        CHECK(log.stops[0].location == "Sunken Arch");
        CHECK(log.stops[0].arriveMinutes == 1479);
        CHECK((int)log.stops[0].items.size() == 2);
        CHECK(log.stops[0].items[0].name == "Gems");
        CHECK(log.stops[0].items[0].quantity == 584633497);
        CHECK(log.stops[0].items[1].name == "Crystals");
        CHECK(log.stops[0].items[1].quantity == 202209422);
        CHECK(log.stops[0].exp == 456423512);
        CHECK(log.stops[0].leaveMinutes == 1562);
        // ── stop 1: Hex 0 Atoll ──
        CHECK((int)log.stops[1].hazards.size() == 1);
        CHECK(log.stops[1].hazards[0] == "Whirlpool");
        CHECK(log.stops[1].location == "Hex 0 Atoll");
        CHECK(log.stops[1].arriveMinutes == 1648);
        CHECK((int)log.stops[1].items.size() == 2);
        CHECK(log.stops[1].items[0].name == "Dark Pearl");
        CHECK(log.stops[1].items[0].quantity == 1175006141);
        CHECK(log.stops[1].items[1].name == "Amber");
        CHECK(log.stops[1].items[1].quantity == 1848890585);
        CHECK(log.stops[1].exp == 415252392);
        CHECK(log.stops[1].leaveMinutes == 1787);
        // ── stop 2: Site 8 Hollow ──
        CHECK((int)log.stops[2].hazards.size() == 1);
        CHECK(log.stops[2].hazards[0] == "Dark Current");
        CHECK(log.stops[2].location == "Site 8 Hollow");
        CHECK(log.stops[2].arriveMinutes == 1831);
        CHECK((int)log.stops[2].items.size() == 2);
        CHECK(log.stops[2].items[0].name == "Resin 4 Block");
        CHECK(log.stops[2].items[0].quantity == 1781911614);
        CHECK(log.stops[2].items[1].name == "Coral");
        CHECK(log.stops[2].items[1].quantity == 1508751711);
        CHECK(log.stops[2].exp == 317796632);
        CHECK(log.stops[2].leaveMinutes == 1990);

        // ── getLocationAt boundaries ──
        CHECK(log.getLocationAt(     0) == "Not yet departed");
        CHECK(log.getLocationAt(  1379) == "Not yet departed");
        CHECK(log.getLocationAt(  1380) == "En route >>>Sunken Arch");
        CHECK(log.getLocationAt(  1429) == "En route >>>Sunken Arch");
        CHECK(log.getLocationAt(  1478) == "En route >>>Sunken Arch");
        CHECK(log.getLocationAt(  1479) == "Sunken Arch");
        CHECK(log.getLocationAt(  1520) == "Sunken Arch");
        CHECK(log.getLocationAt(  1561) == "Sunken Arch");
        CHECK(log.getLocationAt(  1562) == "En route >>>Hex 0 Atoll");
        CHECK(log.getLocationAt(  1605) == "En route >>>Hex 0 Atoll");
        CHECK(log.getLocationAt(  1647) == "En route >>>Hex 0 Atoll");
        CHECK(log.getLocationAt(  1648) == "Hex 0 Atoll");
        CHECK(log.getLocationAt(  1717) == "Hex 0 Atoll");
        CHECK(log.getLocationAt(  1786) == "Hex 0 Atoll");
        CHECK(log.getLocationAt(  1787) == "En route >>>Site 8 Hollow");
        CHECK(log.getLocationAt(  1809) == "En route >>>Site 8 Hollow");
        CHECK(log.getLocationAt(  1830) == "En route >>>Site 8 Hollow");
        CHECK(log.getLocationAt(  1831) == "Site 8 Hollow");
        CHECK(log.getLocationAt(  1910) == "Site 8 Hollow");
        CHECK(log.getLocationAt(  1989) == "Site 8 Hollow");
        CHECK(log.getLocationAt(  1990) == "Returning home");
        CHECK(log.getLocationAt(  2090) == "Returning home");

        // ── getExpAt boundaries ──
        CHECK(log.getExpAt(1478) == 0u);
        CHECK(log.getExpAt(1478) == 0);
        CHECK(log.getExpAt(1479) == 456423512);
        CHECK(log.getExpAt(1647) == 456423512);
        CHECK(log.getExpAt(1648) == 871675904);
        CHECK(log.getExpAt(1830) == 871675904);
        CHECK(log.getExpAt(1831) == 1189472536);
        CHECK(log.getExpAt(1990) == 1189472536);

        // ── encode() outputs pure minutes, not HH:MM ──
        CHECK(
            log.encode() == "Nova Shell[4313]1380?Fog Bank@Sunken Arch:1479$Gems:584633497&Crystals:202209422%456423512:1562?Whirlpool@Hex 0 Atoll:1648$Dark Pearl:1175006141&Amber:1848890585%415252392:1787?Dark Current@Site 8 Hollow:1831$Resin 4 Block:1781911614&Coral:1508751711%317796632:1990");
    }
    {
        auto log = VoyageLog::parse(
            "Nautilus[3332]23?Storm Surge@Storm Cape:24:29$Driftwood:545178500&Saltstone:1087106724%330222870:26?Ice Vortex@Block 2 Shelf:26:55$Crystals:760879933&Ores:159045302%490974549:29:38?Lava Flow@Ember Isle:30:09$Driftwood:221697148&Coconuts:1103580021%376500717:32:42");

        // ── top-level fields ──
        CHECK(log.shipName == "Nautilus");
        CHECK(log.equipment == std::vector<int>({3, 3, 3, 2}));
        CHECK((int)log.stops.size() == 3);
        CHECK(log.departMinutes == 1380);
        // ── stop 0: Storm Cape ──
        CHECK((int)log.stops[0].hazards.size() == 1);
        CHECK(log.stops[0].hazards[0] == "Storm Surge");
        CHECK(log.stops[0].location == "Storm Cape");
        CHECK(log.stops[0].arriveMinutes == 1469);
        CHECK((int)log.stops[0].items.size() == 2);
        CHECK(log.stops[0].items[0].name == "Driftwood");
        CHECK(log.stops[0].items[0].quantity == 545178500);
        CHECK(log.stops[0].items[1].name == "Saltstone");
        CHECK(log.stops[0].items[1].quantity == 1087106724);
        CHECK(log.stops[0].exp == 330222870);
        CHECK(log.stops[0].leaveMinutes == 1560);
        // ── stop 1: Block 2 Shelf ──
        CHECK((int)log.stops[1].hazards.size() == 1);
        CHECK(log.stops[1].hazards[0] == "Ice Vortex");
        CHECK(log.stops[1].location == "Block 2 Shelf");
        CHECK(log.stops[1].arriveMinutes == 1615);
        CHECK((int)log.stops[1].items.size() == 2);
        CHECK(log.stops[1].items[0].name == "Crystals");
        CHECK(log.stops[1].items[0].quantity == 760879933);
        CHECK(log.stops[1].items[1].name == "Ores");
        CHECK(log.stops[1].items[1].quantity == 159045302);
        CHECK(log.stops[1].exp == 490974549);
        CHECK(log.stops[1].leaveMinutes == 1778);
        // ── stop 2: Ember Isle ──
        CHECK((int)log.stops[2].hazards.size() == 1);
        CHECK(log.stops[2].hazards[0] == "Lava Flow");
        CHECK(log.stops[2].location == "Ember Isle");
        CHECK(log.stops[2].arriveMinutes == 1809);
        CHECK((int)log.stops[2].items.size() == 2);
        CHECK(log.stops[2].items[0].name == "Driftwood");
        CHECK(log.stops[2].items[0].quantity == 221697148);
        CHECK(log.stops[2].items[1].name == "Coconuts");
        CHECK(log.stops[2].items[1].quantity == 1103580021);
        CHECK(log.stops[2].exp == 376500717);
        CHECK(log.stops[2].leaveMinutes == 1962);

        // ── getLocationAt boundaries ──
        CHECK(log.getLocationAt(     0) == "Not yet departed");
        CHECK(log.getLocationAt(  1379) == "Not yet departed");
        CHECK(log.getLocationAt(  1380) == "En route >>>Storm Cape");
        CHECK(log.getLocationAt(  1424) == "En route >>>Storm Cape");
        CHECK(log.getLocationAt(  1468) == "En route >>>Storm Cape");
        CHECK(log.getLocationAt(  1469) == "Storm Cape");
        CHECK(log.getLocationAt(  1514) == "Storm Cape");
        CHECK(log.getLocationAt(  1559) == "Storm Cape");
        CHECK(log.getLocationAt(  1560) == "En route >>>Block 2 Shelf");
        CHECK(log.getLocationAt(  1587) == "En route >>>Block 2 Shelf");
        CHECK(log.getLocationAt(  1614) == "En route >>>Block 2 Shelf");
        CHECK(log.getLocationAt(  1615) == "Block 2 Shelf");
        CHECK(log.getLocationAt(  1696) == "Block 2 Shelf");
        CHECK(log.getLocationAt(  1777) == "Block 2 Shelf");
        CHECK(log.getLocationAt(  1778) == "En route >>>Ember Isle");
        CHECK(log.getLocationAt(  1793) == "En route >>>Ember Isle");
        CHECK(log.getLocationAt(  1808) == "En route >>>Ember Isle");
        CHECK(log.getLocationAt(  1809) == "Ember Isle");
        CHECK(log.getLocationAt(  1885) == "Ember Isle");
        CHECK(log.getLocationAt(  1961) == "Ember Isle");
        CHECK(log.getLocationAt(  1962) == "Returning home");
        CHECK(log.getLocationAt(  2062) == "Returning home");

        // ── getExpAt boundaries ──
        CHECK(log.getExpAt(1468) == 0u);
        CHECK(log.getExpAt(1468) == 0);
        CHECK(log.getExpAt(1469) == 330222870);
        CHECK(log.getExpAt(1614) == 330222870);
        CHECK(log.getExpAt(1615) == 821197419);
        CHECK(log.getExpAt(1808) == 821197419);
        CHECK(log.getExpAt(1809) == 1197698136);
        CHECK(log.getExpAt(1962) == 1197698136);

        // ── encode() outputs pure minutes, not HH:MM ──
        CHECK(
            log.encode() == "Nautilus[3332]1380?Storm Surge@Storm Cape:1469$Driftwood:545178500&Saltstone:1087106724%330222870:1560?Ice Vortex@Block 2 Shelf:1615$Crystals:760879933&Ores:159045302%490974549:1778?Lava Flow@Ember Isle:1809$Driftwood:221697148&Coconuts:1103580021%376500717:1962");
    }
}

#endif // _CASE_H_
