#pragma once

#include <functional>
#include <memory>
#include <vector>
#include <unordered_map>

enum class RoomType {
    Casual,
    Ranked,
    Coop,
    Tournament,
    Puzzle
};

class IRoom {
public:
    virtual ~IRoom() = default;
    virtual int GetId() const = 0;
    virtual RoomType GetType() const = 0;
    virtual int GetPlayerCount() const = 0;
};

using RoomPtr = std::shared_ptr<IRoom>;
using RoomCompare = std::function<bool(const RoomPtr&, const RoomPtr&)>;

// Priority rules:
// 1. Rooms whose id appears in idPriority come first, in the order given.
// 2. Other rooms whose type appears in typePriority come next, in the order given.
// 3. All remaining rooms come last.
// 4. Within the same group, higher player count comes first; if still tied, smaller id comes first.
inline RoomCompare BuildRoomCompare(const std::vector<int>& idPriority,
                                    const std::vector<RoomType>& typePriority) {
  std::unordered_map<int, std::size_t> idRank;
  std::unordered_map<RoomType, std::size_t> typeRank;

  for (std::size_t index = 0; index < idPriority.size(); index++) {
    idRank[idPriority[index]] = index;
  }

  for (std::size_t index = 0; index < typePriority.size(); index++) {
    typeRank[typePriority[index]] = index;
  }

  return [idRank = std::move(idRank), typeRank = std::move(typeRank)](
             const RoomPtr& left, const RoomPtr& right) {
    const auto leftId = left->GetId();
    const auto rightId = right->GetId();

    const auto leftIdIter = idRank.find(leftId);
    const auto rightIdIter = idRank.find(rightId);
    const bool leftInIdRank = leftIdIter != idRank.end();
    const bool rightInIdRank = rightIdIter != idRank.end();

    const int leftGroup =
        leftInIdRank ? 0 : (typeRank.count(left->GetType()) ? 1 : 2);
    const int rightGroup =
        rightInIdRank ? 0 : (typeRank.count(right->GetType()) ? 1 : 2);

    if (leftGroup != rightGroup) {
      return leftGroup < rightGroup;
    }

    if (leftGroup == 0 && leftIdIter->second != rightIdIter->second) {
      return leftIdIter->second < rightIdIter->second;
    }

    if (leftGroup == 1) {
      const auto leftTypeRank = typeRank.find(left->GetType())->second;
      const auto rightTypeRank = typeRank.find(right->GetType())->second;
      if (leftTypeRank != rightTypeRank) {
        return leftTypeRank < rightTypeRank;
      }
    }

    if (left->GetPlayerCount() != right->GetPlayerCount()) {
      return left->GetPlayerCount() > right->GetPlayerCount();
    }

    return leftId < rightId;
  };
}
