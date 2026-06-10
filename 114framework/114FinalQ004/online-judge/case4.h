#ifndef _CASE_H_
#define _CASE_H_

#define SECRET "936e1dadf18e7f72d4340efa8a8affc4b93ca705336c2afd8bae821267a08f2d"

#include <algorithm>

#include "room_test_utils.h"
#include "test.h"

TEST_CASE("Rooms outside all priorities fall back to player count") {
    auto compare = BuildRoomCompare({100}, {RoomType::Ranked});

    std::vector<RoomPtr> rooms = {
        MakeRoom(10, RoomType::Casual, 5),
        MakeRoom(11, RoomType::Puzzle, 20),
        MakeRoom(12, RoomType::Puzzle, 5),
        MakeRoom(100, RoomType::Casual, 1),
        MakeRoom(13, RoomType::Ranked, 7),
    };

    std::sort(rooms.begin(), rooms.end(), compare);

    CHECK(rooms[0]->GetId() == 100);
    CHECK(rooms[1]->GetId() == 13);
    CHECK(rooms[2]->GetId() == 11);
    CHECK(rooms[3]->GetId() == 10);
    CHECK(rooms[4]->GetId() == 12);
}

#endif
