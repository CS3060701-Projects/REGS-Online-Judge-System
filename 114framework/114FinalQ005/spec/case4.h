#ifndef _CASE_H_
#define _CASE_H_
#define SECRET "@SECRET@"

#include "to_base.h"
#include "test.h"

TEST_CASE("case: 8,2") {
  std::array<uint8_t, 8> buf{0};
  auto result = to_base<2>(buf);

  auto answer = "0";

  REQUIRE(result == answer);
}

#endif