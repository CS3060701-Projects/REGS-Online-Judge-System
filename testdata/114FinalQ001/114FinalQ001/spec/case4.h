#ifndef _CASE_H_
#define _CASE_H_
#define SECRET "@SECRET@"
#include "media_type.h"
#include "test.h"

namespace filter {
inline Pipe<H264Frame, H264Frame> myfilter() {
  return Pipe<H264Frame, H264Frame>(
      [](H264Frame in, H264Frame& out) {
        out.data[3] = 0x00;
        out.data[9] = 0x15;
        return true;
      });
}
}  // namespace filter

TEST_CASE("case: custom filter corruption", "[pipeline][filter]") {
  auto source = initSource(0x98);
  RawChunk result;

  // clang-format off
  auto pipeline = codecs::h264enc()
    | filter::myfilter()
    | codecs::h264dec();
  // clang-format on

  bool success = pipeline.process(source, result);

  CHECK(success == true);
  CHECK(result.data[3] == (0x00 ^ 0x65));
  CHECK(result.data[9] == (0x15 ^ 0x65));
}

#endif