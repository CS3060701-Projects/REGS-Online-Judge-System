#ifndef _TO_BASE_H_
#define _TO_BASE_H_
#include <algorithm>
#include <array>
#include <cassert>
#include <cstdint>
#include <string>
#include <vector>

#include "generate.h"

namespace {
const char charset[] = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ";
}

template <std::size_t N>
inline std::string debug(std::array<uint8_t, N>& buf) {
  static_assert(N % 8 == 0, "N must be a multiple of 8");

  // Map uint8_t[N] to uint32_t[N]
  using chunk_type = uint32_t;
  constexpr size_t chunkSize = N / sizeof(chunk_type);
  std::array<chunk_type, chunkSize> chunks{0};
  for (size_t i = 0; i < chunkSize; ++i) {
    chunks[i] = (static_cast<uint32_t>(buf[i * 4 + 0]) << 24) |
                (static_cast<uint32_t>(buf[i * 4 + 1]) << 16) |
                (static_cast<uint32_t>(buf[i * 4 + 2]) << 8) |
                (static_cast<uint32_t>(buf[i * 4 + 3]));
    std::printf("chunkSize[%d] = %08x \n", i, chunks[i]);
  }
}

template <int base, std::size_t N>
inline std::string to_base(std::array<uint8_t, N>& buf) {
  static_assert(2 <= base && base <= 36, "base must be between 2 and 36");
  static_assert(N % 8 == 0, "N must be a multiple of 8");

  // Map uint8_t[N] to uint32_t[N]
  using chunk_type = uint32_t;
  constexpr size_t chunkSize = N / sizeof(chunk_type);
  std::vector<chunk_type> chunks(chunkSize, 0);
  for (size_t i = 0; i < chunkSize; ++i) {
    chunks[i] = (static_cast<uint32_t>(buf[i * 4 + 0]) << 24) |
                (static_cast<uint32_t>(buf[i * 4 + 1]) << 16) |
                (static_cast<uint32_t>(buf[i * 4 + 2]) << 8) |
                (static_cast<uint32_t>(buf[i * 4 + 3]));
  }

  std::string result = "";

  size_t start = 0;
  while (start < chunkSize && chunks[start] == 0) ++start;
  if (start == chunkSize) return "0";

  auto done = [&]() {
    return std::all_of(chunks.begin(), chunks.end(),
                [](auto chunk) { return chunk == 0; });
  };
  while (!done()) {
    while (chunks[start] == 0) ++start;
    uint64_t remainder = 0;

    for (size_t i = 0; i < chunkSize; ++i) {
      uint64_t current = (remainder << 32) | chunks[i];
      chunks[i] = static_cast<chunk_type>(current / base);
      remainder = static_cast<chunk_type>(current % base);
    }

    result += charset[remainder];
  }

  std::reverse(result.begin(), result.end());
  return result;
}

#endif