#ifndef _CASE_H_
#define _CASE_H_
#define SECRET "@SECRET@"
#include "media_type.h"
#include "test.h"

TEST_CASE("case: vp9-MKV convert", "[pipeline][vp9][mkv]") {
  auto source = initSource(0x66);

  // clang-format off
  MKVPacket middle;
  auto pipeline = codecs::vp9enc() | muldex::vp9ToMKV();

  RawChunk result;
  auto pipeline2 = muldex::MKVTovp9()| codecs::vp9dec();
  // clang-format on

  bool checkpoint1 = pipeline.process(source, middle);
  CHECK(checkpoint1 == true);

  bool checkpoint2 = pipeline2.process(middle, result);
  CHECK(checkpoint2 == true);
  CHECK(source == result);
}

TEST_CASE("case: h265-MKV convert", "[pipeline][h265][mkv]") {
  auto source = initSource(0x15);
  RawChunk result;

  // clang-format off
  auto pipeline = codecs::h265enc() 
    | muldex::h265ToMKV()
    | muldex::MKVToh265()
    | codecs::h265dec();
  // clang-format on

  bool success = pipeline.process(source, result);
  CHECK(success == true);
  CHECK(source == result);
}

#endif