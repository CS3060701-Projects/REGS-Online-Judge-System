#ifndef _CASE_H_
#define _CASE_H_

#define SECRET "@SECRET@"

#include <algorithm>

#include "room_test_utils.h"
#include "test.h"

TEST_CASE("Type priority orders the middle layer") {
  auto compare = BuildRoomCompare({}, {RoomType::Tournament, RoomType::Ranked});

  std::vector<RoomPtr> rooms = {
      MakeRoom(1, RoomType::Ranked, 12),     MakeRoom(2, RoomType::Ranked, 40),
      MakeRoom(3, RoomType::Tournament, 20), MakeRoom(4, RoomType::Casual, 99),
      MakeRoom(5, RoomType::Casual, 2),
  };

  std::sort(rooms.begin(), rooms.end(), compare);

  CHECK(rooms[0]->GetId() == 3);
  CHECK(rooms[1]->GetId() == 2);
  CHECK(rooms[2]->GetId() == 1);
  CHECK(rooms[3]->GetId() == 4);
  CHECK(rooms[4]->GetId() == 5);
}

#endif