#pragma once

#include "game_room_sort.h"

class TestRoom : public IRoom {
public:
    TestRoom(int id, RoomType type, int playerCount)
        : id_(id), type_(type), playerCount_(playerCount) {}

    int GetId() const override {
        return id_;
    }

    RoomType GetType() const override {
        return type_;
    }

    int GetPlayerCount() const override {
        return playerCount_;
    }

private:
    int id_;
    RoomType type_;
    int playerCount_;
};

inline RoomPtr MakeRoom(int id, RoomType type, int playerCount) {
    return std::make_shared<TestRoom>(id, type, playerCount);
}
