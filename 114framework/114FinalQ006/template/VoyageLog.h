#ifndef VOYAGE_LOG_H_
#define VOYAGE_LOG_H_

#include <string>
#include <vector>
#include <sstream>
#include <iomanip>

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

    static VoyageLog parse(const std::string& input);

    std::string getLocationAt(int minutes) const;

    long long getExpAt(int minutes) const;

    std::string encode() const;
    
    std::string display() const{
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
    };

private:
    static std::string fmt(int minutes) {
        std::ostringstream oss;
        oss << ">-- Elapsed " << minutes / 60 << "h " << std::setw(2) << std::setfill('0') << minutes % 60 << "m --<\n";
        return oss.str();
    }

    
};

#endif // VOYAGE_LOG_H_
