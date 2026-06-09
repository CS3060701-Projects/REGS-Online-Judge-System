#ifndef _CASE_H_
#define _CASE_H_
#define SECRET "@SECRET@"

#include "VoyageLog.h"
#include "test.h"

TEST_CASE("Parse - single stop") {
    {
        // Nova Shell[2341]7:15@Pearl Ridge:8:21$Coral:43%14352:10:55
        auto log = VoyageLog::parse("Nova Shell[2341]7:15@Pearl Ridge:8:21$Coral:43%14352:10:55");

        CHECK(log.shipName == "Nova Shell");
        CHECK(log.equipment == std::vector<int>({2, 3, 4, 1}));
        CHECK((int)log.stops.size() == 1);
        CHECK(log.departMinutes == 435);

        CHECK(log.stops[0].location == "Pearl Ridge");
        CHECK(log.stops[0].arriveMinutes == 501);
        CHECK((int)log.stops[0].hazards.size() == 0);
        CHECK((int)log.stops[0].items.size() == 1);
        CHECK(log.stops[0].items[0].name == "Coral");
        CHECK(log.stops[0].items[0].quantity == 43);
        CHECK(log.stops[0].exp == 14352);
        CHECK(log.stops[0].leaveMinutes == 655);
    }
    {
        // Subber[3221]18:30@Jade Grotto:19:24$Sea Glass:31%13235:21:54
        auto log = VoyageLog::parse("Subber[3221]18:30@Jade Grotto:19:24$Sea Glass:31%13235:21:54");

        CHECK(log.shipName == "Subber");
        CHECK(log.equipment == std::vector<int>({3, 2, 2, 1}));
        CHECK((int)log.stops.size() == 1);
        CHECK(log.departMinutes == 1110);

        CHECK(log.stops[0].location == "Jade Grotto");
        CHECK(log.stops[0].arriveMinutes == 1164);
        CHECK((int)log.stops[0].hazards.size() == 0);
        CHECK((int)log.stops[0].items.size() == 1);
        CHECK(log.stops[0].items[0].name == "Sea Glass");
        CHECK(log.stops[0].items[0].quantity == 31);
        CHECK(log.stops[0].exp == 13235);
        CHECK(log.stops[0].leaveMinutes == 1314);
    }
    {
        // Blue Phantom[2144]7@Silver Lagoon:8:53$Saltstone:51%12429:10:47
        auto log = VoyageLog::parse("Blue Phantom[2144]7@Silver Lagoon:8:53$Saltstone:51%12429:10:47");

        CHECK(log.shipName == "Blue Phantom");
        CHECK(log.equipment == std::vector<int>({2, 1, 4, 4}));
        CHECK((int)log.stops.size() == 1);
        CHECK(log.departMinutes == 420);

        CHECK(log.stops[0].location == "Silver Lagoon");
        CHECK(log.stops[0].arriveMinutes == 533);
        CHECK((int)log.stops[0].hazards.size() == 0);
        CHECK((int)log.stops[0].items.size() == 1);
        CHECK(log.stops[0].items[0].name == "Saltstone");
        CHECK(log.stops[0].items[0].quantity == 51);
        CHECK(log.stops[0].exp == 12429);
        CHECK(log.stops[0].leaveMinutes == 647);
    }
}

#endif // _CASE_H_
