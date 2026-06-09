#ifndef _CASE_H_
#define _CASE_H_

#define SECRET "d5956684846393dcefb01af9135a46e1065e3cc770f1b6610bddd28dd8fe60bc"

#include <algorithm>

#include "room_test_utils.h"
#include "test.h"

TEST_CASE("ID priority beats type priority") {
    auto compare = BuildRoomCompare({9}, {RoomType::Coop});

    std::vector<RoomPtr> rooms = {
        MakeRoom(2, RoomType::Casual, 50),
        MakeRoom(1, RoomType::Coop, 100),
        MakeRoom(9, RoomType::Casual, 1),
    };

    std::sort(rooms.begin(), rooms.end(), compare);

    CHECK(rooms[0]->GetId() == 9);
    CHECK(rooms[1]->GetId() == 1);
    CHECK(rooms[2]->GetId() == 2);
}

#endif
