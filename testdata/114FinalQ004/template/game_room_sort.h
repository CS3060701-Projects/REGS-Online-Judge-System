#pragma once

#include <functional>
#include <memory>
#include <vector>

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
RoomCompare BuildRoomCompare(const std::vector<int>& idPriority,
                             const std::vector<RoomType>& typePriority);
