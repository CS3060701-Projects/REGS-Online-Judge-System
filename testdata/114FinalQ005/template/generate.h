#ifndef _GENERATE_H_
#define _GENERATE_H_

#include <array>
#include <cstdint>
#include <iomanip>
#include <iostream>

namespace gen {

class RandEngine {
  uint32_t state;

 public:
  explicit RandEngine(uint32_t seed) : state(seed == 0 ? 0xbadcafe : seed) {}

  uint8_t next() {
    state ^= (state << 13);
    state ^= (state >> 17);
    state ^= (state << 5);
    return static_cast<uint8_t>(state & 0xFF);
  }
};

template <std::size_t N>
inline void fill(std::array<uint8_t, N>& arr, uint32_t seed) {
  RandEngine rng(seed);
  for (std::size_t i = 0; i < N; ++i) {
    arr[i] = rng.next();
  }
}

template <typename T, std::size_t N>
inline void toHexString(std::array<T, N>& arr, std::string delim = "") {
  auto originFlags = std::cout.flags();
  std::cout << std::hex << std::setfill('0') << "0x";
  for (auto byte : arr) {
    std::cout << std::setw(2) << static_cast<int>(byte) << delim;
  }
  std::cout << 'n' << std::endl;
  std::cout.flags(originFlags);
}
}  // namespace gen
#endif