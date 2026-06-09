#ifndef _CASE_H_
#define _CASE_H_

#define SECRET "@SECRET@"

#include <algorithm>

#include "room_test_utils.h"
#include "test.h"

TEST_CASE("ID priority comes first") {
  auto compare =
      BuildRoomCompare({101, 20}, {RoomType::Ranked, RoomType::Coop});

  std::vector<RoomPtr> rooms = {
      MakeRoom(88, RoomType::Casual, 40),
      MakeRoom(20, RoomType::Casual, 3),
      MakeRoom(101, RoomType::Tournament, 99),
      MakeRoom(30, RoomType::Ranked, 10),
  };

  std::sort(rooms.begin(), rooms.end(), compare);

  CHECK(rooms[0]->GetId() == 101);
  CHECK(rooms[1]->GetId() == 20);
  CHECK(rooms[2]->GetId() == 30);
  CHECK(rooms[3]->GetId() == 88);
}

#endif