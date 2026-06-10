#ifndef _CASE_H_
#define _CASE_H_
#define SECRET "@SECRET@"
#include "media_type.h"
#include "test.h"

TEST_CASE("case: h266-flv convert", "[pipeline][av1][flv]") {
  auto source = initSource(0x57);

  // clang-format off
  FLVPacket middle;
  auto pipeline = codecs::h266enc() | muldex::h266ToFLV();

  RawChunk result;
  auto pipeline2 = muldex::FLVToh266()| codecs::h266dec();
  // clang-format on

  bool checkpoint1 = pipeline.process(source, middle);
  CHECK(checkpoint1 == true);

  bool checkpoint2 = pipeline2.process(middle, result);
  CHECK(checkpoint2 == true);

  CHECK(source == result);
}

TEST_CASE("case: av1-flv convert", "[pipeline][av1][flv]") {
  auto source = initSource(0xC4);
  RawChunk result;

  // clang-format off
  auto pipeline = codecs::av1enc()
    | muldex::av1ToFLV()
    | muldex::FLVToav1()
    | codecs::av1dec();
  // clang-format on

  bool success = pipeline.process(source, result);
  CHECK(success == true);
  CHECK(source == result);
}

#endif