#ifndef _CASE_H_
#define _CASE_H_

#define SECRET "2a101ead67ccb9b5c4bc791415b7cd0c618766c8b0b4a5014c75099c451941df"

#include <algorithm>

#include "room_test_utils.h"
#include "test.h"

TEST_CASE("Mixed priorities produce one total order") {
    auto compare = BuildRoomCompare({42, 7}, {RoomType::Tournament, RoomType::Ranked, RoomType::Coop});

    std::vector<RoomPtr> rooms = {
        MakeRoom(3, RoomType::Casual, 40),
        MakeRoom(1, RoomType::Coop, 30),
        MakeRoom(42, RoomType::Casual, 8),
        MakeRoom(8, RoomType::Tournament, 20),
        MakeRoom(7, RoomType::Puzzle, 15),
        MakeRoom(4, RoomType::Casual, 40),
        MakeRoom(9, RoomType::Ranked, 30),
        MakeRoom(5, RoomType::Tournament, 60),
    };

    std::sort(rooms.begin(), rooms.end(), compare);

    CHECK(rooms[0]->GetId() == 42);
    CHECK(rooms[1]->GetId() == 7);
    CHECK(rooms[2]->GetId() == 5);
    CHECK(rooms[3]->GetId() == 8);
    CHECK(rooms[4]->GetId() == 9);
    CHECK(rooms[5]->GetId() == 1);
    CHECK(rooms[6]->GetId() == 3);
    CHECK(rooms[7]->GetId() == 4);
}

#endif
