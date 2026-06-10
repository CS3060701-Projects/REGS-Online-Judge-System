#ifndef VOYAGE_LOG_H_
#define VOYAGE_LOG_H_

#include <cctype>
#include <iomanip>
#include <sstream>
#include <string>
#include <vector>

struct Item {
    std::string name;
    int quantity;
};

struct Stop {
    std::vector<std::string> hazards;
    std::string location;
    int arriveMinutes;
    std::vector<Item> items;
    int exp;
    int leaveMinutes;
};

class VoyageLog {
public:
    std::string shipName;
    std::vector<int> equipment;
    int departMinutes;
    std::vector<Stop> stops;

    static VoyageLog parse(const std::string& s) {
        VoyageLog log;

        std::istringstream iss(s + "\x1F");
        std::stringstream buf("");
        char cursor;

        enum class State { SHIP_NAME, EQUIPMENT, TIME, HAZARD, LOCATION, ITEM_NAME, ITEM_QTY, EXP };
        enum class TimeRole { DEPART, ARRIVE, LEAVE };

        State state = State::SHIP_NAME;
        TimeRole timeRole = TimeRole::DEPART;
        bool inMM = false;
        int timeH = 0;

        Stop currentStop;
        Item currentItem;

        auto take = [&]() -> std::string {
            auto r = buf.str();
            buf.str("");
            buf.clear();
            return r;
        };

        while (iss.get(cursor)) {
            switch (state) {
            case State::SHIP_NAME:
                if (cursor == '[') {
                    log.shipName = take();
                    state = State::EQUIPMENT;
                } else {
                    buf << cursor;
                }
                break;
            case State::EQUIPMENT:
                if (cursor == ']') {
                    state = State::TIME;
                    timeRole = TimeRole::DEPART;
                } else {
                    log.equipment.push_back(cursor - '0');
                }
                break;
            case State::TIME:
                if (isdigit(cursor)) {
                    buf << cursor;
                } else if (cursor == ':') {
                    timeH = std::stoi(take());
                    inMM = true;
                } else {
                    std::string raw = take();
                    int minutes = inMM ? timeH * 60 + std::stoi(raw) : std::stoi(raw) * 60;
                    inMM = false;
                    timeH = 0;

                    switch (timeRole) {
                    case TimeRole::DEPART:
                        log.departMinutes = minutes;
                        currentStop = Stop{};
                        state = (cursor == '?') ? State::HAZARD : State::LOCATION;
                        break;

                    case TimeRole::ARRIVE:
                        currentStop.arriveMinutes = minutes;
                        state = State::ITEM_NAME;
                        break;

                    case TimeRole::LEAVE:
                        currentStop.leaveMinutes = minutes;
                        log.stops.push_back(currentStop);
                        if (cursor != '\x1F') {
                            currentStop = Stop{};
                            state = (cursor == '?') ? State::HAZARD : State::LOCATION;
                        }
                        break;
                    }
                }
                break;
            case State::HAZARD:
                if (cursor == '@' || cursor == '?') {
                    auto hazard = take();
                    if (!hazard.empty())
                        currentStop.hazards.push_back(hazard);
                    if (cursor == '@')
                        state = State::LOCATION;
                } else {
                    buf << cursor;
                }
                break;
            case State::LOCATION:
                if (cursor == ':') {
                    currentStop.location = take();
                    state = State::TIME;
                    timeRole = TimeRole::ARRIVE;
                } else {
                    buf << cursor;
                }
                break;
            case State::ITEM_NAME:
                if (cursor == ':') {
                    currentItem = Item{take(), 0};
                    state = State::ITEM_QTY;
                } else {
                    buf << cursor;
                }
                break;
            case State::ITEM_QTY:
                if (cursor == '&' || cursor == '%' || cursor == ':') {
                    currentItem.quantity = std::stoi(take());
                    currentStop.items.push_back(currentItem);
                    if (cursor == '&')
                        state = State::ITEM_NAME;
                    else if (cursor == '%')
                        state = State::EXP;
                    else {
                        state = State::TIME;
                        timeRole = TimeRole::LEAVE;
                    }
                } else {
                    buf << cursor;
                }
                break;
            case State::EXP:
                if (cursor == ':') {
                    currentStop.exp = std::stoi(take());
                    state = State::TIME;
                    timeRole = TimeRole::LEAVE;
                } else {
                    buf << cursor;
                }
                break;
            }
        }

        return log;
    }

    std::string getLocationAt(int minutes) const {
        if (stops.empty() || minutes < departMinutes)
            return "Not yet departed";

        for (size_t k = 0; k < stops.size(); ++k) {
            int enRouteStart = (k == 0) ? departMinutes : stops[k - 1].leaveMinutes;
            if (minutes >= enRouteStart && minutes < stops[k].arriveMinutes)
                return "En route >>>" + stops[k].location;

            if (minutes >= stops[k].arriveMinutes && minutes < stops[k].leaveMinutes)
                return stops[k].location;
        }

        return "Returning home";
    }

    long long getExpAt(int minutes) const {
        long long total = 0;
        for (const auto& stop : stops)
            if (stop.arriveMinutes <= minutes)
                total += stop.exp;
        return total;
    }

    std::string display() const {
        static const std::vector<std::string> parts = {"Bow", "Hull", "Stern", "Bridge"};
        std::ostringstream oss;

        oss << shipName << " sets sail\nEquipment: ";
        for (size_t k = 0; k < equipment.size(); ++k) {
            if (k > 0)
                oss << ", ";
            oss << "Grade-" << equipment[k] << " " << parts[k];
        }
        oss << "\n---\n";
        oss << fmt(departMinutes);

        for (const auto& stop : stops) {
            for (const auto& h : stop.hazards)
                oss << "Encountered hazard: " << h << "\n";
            oss << "Heading to " << stop.location << "\n";
            oss << fmt(stop.arriveMinutes);
            oss << "Arrived at " << stop.location << "\n";
            for (const auto& item : stop.items)
                oss << "Obtained " << item.name << " x" << item.quantity << "\n";
            oss << stop.location << " exploration complete\n";
            oss << shipName << " gained " << stop.exp << " EXP\n";
            oss << fmt(stop.leaveMinutes);
        }

        oss << "All locations explored. Returning home.\n";
        return oss.str();
    }

    std::string encode() const {
        std::ostringstream oss;

        oss << shipName << "[";
        for (int e : equipment)
            oss << e;
        oss << "]" << departMinutes;

        for (const auto& s : stops) {
            for (const auto& h : s.hazards)
                oss << "?" << h;

            oss << "@" << s.location << ":" << s.arriveMinutes;

            for (size_t j = 0; j < s.items.size(); ++j)
                oss << (j == 0 ? "$" : "&") << s.items[j].name << ":" << s.items[j].quantity;

            oss << "%" << s.exp << ":" << s.leaveMinutes;
        }

        return oss.str();
    }

private:
    static std::string fmt(int minutes) {
        std::ostringstream oss;
        oss << ">-- Elapsed " << minutes / 60 << "h " << std::setw(2) << std::setfill('0') << minutes % 60 << "m --<\n";
        return oss.str();
    }
};

#endif // VOYAGE_LOG_H_
