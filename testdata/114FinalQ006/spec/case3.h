#ifndef _CASE_H_
#define _CASE_H_
#define SECRET "@SECRET@"

#include "VoyageLog.h"
#include "test.h"

TEST_CASE("getLocationAt - 3-stop boundary sweep") {
    {
        auto log = VoyageLog::parse(
            "Nautilus[2213]6:45@Misty Inlet:8:15$Gems:263%34639:11:14@Pearl Ridge:12:53$Dark Pearl:842%7960:13:31@Twin Peaks:15:29$Ores:505%19000:17:39");

        // ── top-level fields ──
        CHECK(log.shipName == "Nautilus");
        CHECK(log.equipment == std::vector<int>({2, 2, 1, 3}));
        CHECK((int)log.stops.size() == 3);
        CHECK(log.departMinutes == 405);
        // ── stop 0: Misty Inlet ──
        CHECK(log.stops[0].location == "Misty Inlet");
        CHECK(log.stops[0].arriveMinutes == 495);
        CHECK((int)log.stops[0].hazards.size() == 0);
        CHECK((int)log.stops[0].items.size() == 1);
        CHECK(log.stops[0].items[0].name == "Gems");
        CHECK(log.stops[0].items[0].quantity == 263);
        CHECK(log.stops[0].exp == 34639);
        CHECK(log.stops[0].leaveMinutes == 674);
        // ── stop 1: Pearl Ridge ──
        CHECK(log.stops[1].location == "Pearl Ridge");
        CHECK(log.stops[1].arriveMinutes == 773);
        CHECK((int)log.stops[1].hazards.size() == 0);
        CHECK((int)log.stops[1].items.size() == 1);
        CHECK(log.stops[1].items[0].name == "Dark Pearl");
        CHECK(log.stops[1].items[0].quantity == 842);
        CHECK(log.stops[1].exp == 7960);
        CHECK(log.stops[1].leaveMinutes == 811);
        // ── stop 2: Twin Peaks ──
        CHECK(log.stops[2].location == "Twin Peaks");
        CHECK(log.stops[2].arriveMinutes == 929);
        CHECK((int)log.stops[2].hazards.size() == 0);
        CHECK((int)log.stops[2].items.size() == 1);
        CHECK(log.stops[2].items[0].name == "Ores");
        CHECK(log.stops[2].items[0].quantity == 505);
        CHECK(log.stops[2].exp == 19000);
        CHECK(log.stops[2].leaveMinutes == 1059);

        // ── getLocationAt boundaries ──
        CHECK(log.getLocationAt(     0) == "Not yet departed");
        CHECK(log.getLocationAt(   404) == "Not yet departed");
        CHECK(log.getLocationAt(   405) == "En route >>>Misty Inlet");
        CHECK(log.getLocationAt(   450) == "En route >>>Misty Inlet");
        CHECK(log.getLocationAt(   494) == "En route >>>Misty Inlet");
        CHECK(log.getLocationAt(   495) == "Misty Inlet");
        CHECK(log.getLocationAt(   584) == "Misty Inlet");
        CHECK(log.getLocationAt(   673) == "Misty Inlet");
        CHECK(log.getLocationAt(   674) == "En route >>>Pearl Ridge");
        CHECK(log.getLocationAt(   723) == "En route >>>Pearl Ridge");
        CHECK(log.getLocationAt(   772) == "En route >>>Pearl Ridge");
        CHECK(log.getLocationAt(   773) == "Pearl Ridge");
        CHECK(log.getLocationAt(   792) == "Pearl Ridge");
        CHECK(log.getLocationAt(   810) == "Pearl Ridge");
        CHECK(log.getLocationAt(   811) == "En route >>>Twin Peaks");
        CHECK(log.getLocationAt(   870) == "En route >>>Twin Peaks");
        CHECK(log.getLocationAt(   928) == "En route >>>Twin Peaks");
        CHECK(log.getLocationAt(   929) == "Twin Peaks");
        CHECK(log.getLocationAt(   994) == "Twin Peaks");
        CHECK(log.getLocationAt(  1058) == "Twin Peaks");
        CHECK(log.getLocationAt(  1059) == "Returning home");
        CHECK(log.getLocationAt(  1159) == "Returning home");
    }
    {
        auto log = VoyageLog::parse("Arctic Fin[1431]5:15@Frozen Fjord:6:11$Kelp Fiber:815%7063:9:11@Ember Isle:11:10$Frost Shard:669%16482:12:43@Sunken Arch:14:20$Pearls:905%17846:15:53");

        // ── top-level fields ──
        CHECK(log.shipName == "Arctic Fin");
        CHECK(log.equipment == std::vector<int>({1, 4, 3, 1}));
        CHECK((int)log.stops.size() == 3);
        CHECK(log.departMinutes == 315);
        // ── stop 0: Frozen Fjord ──
        CHECK(log.stops[0].location == "Frozen Fjord");
        CHECK(log.stops[0].arriveMinutes == 371);
        CHECK((int)log.stops[0].hazards.size() == 0);
        CHECK((int)log.stops[0].items.size() == 1);
        CHECK(log.stops[0].items[0].name == "Kelp Fiber");
        CHECK(log.stops[0].items[0].quantity == 815);
        CHECK(log.stops[0].exp == 7063);
        CHECK(log.stops[0].leaveMinutes == 551);
        // ── stop 1: Ember Isle ──
        CHECK(log.stops[1].location == "Ember Isle");
        CHECK(log.stops[1].arriveMinutes == 670);
        CHECK((int)log.stops[1].hazards.size() == 0);
        CHECK((int)log.stops[1].items.size() == 1);
        CHECK(log.stops[1].items[0].name == "Frost Shard");
        CHECK(log.stops[1].items[0].quantity == 669);
        CHECK(log.stops[1].exp == 16482);
        CHECK(log.stops[1].leaveMinutes == 763);
        // ── stop 2: Sunken Arch ──
        CHECK(log.stops[2].location == "Sunken Arch");
        CHECK(log.stops[2].arriveMinutes == 860);
        CHECK((int)log.stops[2].hazards.size() == 0);
        CHECK((int)log.stops[2].items.size() == 1);
        CHECK(log.stops[2].items[0].name == "Pearls");
        CHECK(log.stops[2].items[0].quantity == 905);
        CHECK(log.stops[2].exp == 17846);
        CHECK(log.stops[2].leaveMinutes == 953);

        // ── getLocationAt boundaries ──
        CHECK(log.getLocationAt(     0) == "Not yet departed");
        CHECK(log.getLocationAt(   314) == "Not yet departed");
        CHECK(log.getLocationAt(   315) == "En route >>>Frozen Fjord");
        CHECK(log.getLocationAt(   343) == "En route >>>Frozen Fjord");
        CHECK(log.getLocationAt(   370) == "En route >>>Frozen Fjord");
        CHECK(log.getLocationAt(   371) == "Frozen Fjord");
        CHECK(log.getLocationAt(   461) == "Frozen Fjord");
        CHECK(log.getLocationAt(   550) == "Frozen Fjord");
        CHECK(log.getLocationAt(   551) == "En route >>>Ember Isle");
        CHECK(log.getLocationAt(   610) == "En route >>>Ember Isle");
        CHECK(log.getLocationAt(   669) == "En route >>>Ember Isle");
        CHECK(log.getLocationAt(   670) == "Ember Isle");
        CHECK(log.getLocationAt(   716) == "Ember Isle");
        CHECK(log.getLocationAt(   762) == "Ember Isle");
        CHECK(log.getLocationAt(   763) == "En route >>>Sunken Arch");
        CHECK(log.getLocationAt(   811) == "En route >>>Sunken Arch");
        CHECK(log.getLocationAt(   859) == "En route >>>Sunken Arch");
        CHECK(log.getLocationAt(   860) == "Sunken Arch");
        CHECK(log.getLocationAt(   906) == "Sunken Arch");
        CHECK(log.getLocationAt(   952) == "Sunken Arch");
        CHECK(log.getLocationAt(   953) == "Returning home");
        CHECK(log.getLocationAt(  1053) == "Returning home");
    }
    {
        auto log = VoyageLog::parse("Arctic Fin[3231]18:15@Lava Shelf:19$Magma Rock:96%11973:20:14@Shadow Reef:21:32$Sea Glass:970%37697:22:11@Amber Shoal:24:02$Amber:987%46986:25:32");

        // ── top-level fields ──
        CHECK(log.shipName == "Arctic Fin");
        CHECK(log.equipment == std::vector<int>({3, 2, 3, 1}));
        CHECK((int)log.stops.size() == 3);
        CHECK(log.departMinutes == 1095);
        // ── stop 0: Lava Shelf ──
        CHECK(log.stops[0].location == "Lava Shelf");
        CHECK(log.stops[0].arriveMinutes == 1140);
        CHECK((int)log.stops[0].hazards.size() == 0);
        CHECK((int)log.stops[0].items.size() == 1);
        CHECK(log.stops[0].items[0].name == "Magma Rock");
        CHECK(log.stops[0].items[0].quantity == 96);
        CHECK(log.stops[0].exp == 11973);
        CHECK(log.stops[0].leaveMinutes == 1214);
        // ── stop 1: Shadow Reef ──
        CHECK(log.stops[1].location == "Shadow Reef");
        CHECK(log.stops[1].arriveMinutes == 1292);
        CHECK((int)log.stops[1].hazards.size() == 0);
        CHECK((int)log.stops[1].items.size() == 1);
        CHECK(log.stops[1].items[0].name == "Sea Glass");
        CHECK(log.stops[1].items[0].quantity == 970);
        CHECK(log.stops[1].exp == 37697);
        CHECK(log.stops[1].leaveMinutes == 1331);
        // ── stop 2: Amber Shoal ──
        CHECK(log.stops[2].location == "Amber Shoal");
        CHECK(log.stops[2].arriveMinutes == 1442);
        CHECK((int)log.stops[2].hazards.size() == 0);
        CHECK((int)log.stops[2].items.size() == 1);
        CHECK(log.stops[2].items[0].name == "Amber");
        CHECK(log.stops[2].items[0].quantity == 987);
        CHECK(log.stops[2].exp == 46986);
        CHECK(log.stops[2].leaveMinutes == 1532);

        // ── getLocationAt boundaries ──
        CHECK(log.getLocationAt(     0) == "Not yet departed");
        CHECK(log.getLocationAt(  1094) == "Not yet departed");
        CHECK(log.getLocationAt(  1095) == "En route >>>Lava Shelf");
        CHECK(log.getLocationAt(  1117) == "En route >>>Lava Shelf");
        CHECK(log.getLocationAt(  1139) == "En route >>>Lava Shelf");
        CHECK(log.getLocationAt(  1140) == "Lava Shelf");
        CHECK(log.getLocationAt(  1177) == "Lava Shelf");
        CHECK(log.getLocationAt(  1213) == "Lava Shelf");
        CHECK(log.getLocationAt(  1214) == "En route >>>Shadow Reef");
        CHECK(log.getLocationAt(  1253) == "En route >>>Shadow Reef");
        CHECK(log.getLocationAt(  1291) == "En route >>>Shadow Reef");
        CHECK(log.getLocationAt(  1292) == "Shadow Reef");
        CHECK(log.getLocationAt(  1311) == "Shadow Reef");
        CHECK(log.getLocationAt(  1330) == "Shadow Reef");
        CHECK(log.getLocationAt(  1331) == "En route >>>Amber Shoal");
        CHECK(log.getLocationAt(  1386) == "En route >>>Amber Shoal");
        CHECK(log.getLocationAt(  1441) == "En route >>>Amber Shoal");
        CHECK(log.getLocationAt(  1442) == "Amber Shoal");
        CHECK(log.getLocationAt(  1487) == "Amber Shoal");
        CHECK(log.getLocationAt(  1531) == "Amber Shoal");
        CHECK(log.getLocationAt(  1532) == "Returning home");
        CHECK(log.getLocationAt(  1632) == "Returning home");
    }
}

#endif // _CASE_H_
