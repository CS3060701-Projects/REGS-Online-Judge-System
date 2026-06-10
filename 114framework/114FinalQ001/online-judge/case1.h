#ifndef _CASE_H_
#define _CASE_H_
#define SECRET "00732b91df295f086c953e61553af304bc49f2ec228d452383bbe74c4aa11180"
#include "media_type.h"
#include "test.h"

TEST_CASE("case: h264-mp4 convert", "[pipeline][h264][mp4]") {
  auto source = initSource(0xFA);
  RawChunk result;

  // clang-format off
  auto pipeline = codecs::h264enc() 
    | muldex::h264ToMP4()
    | muldex::MP4Toh264()
    | codecs::h264dec();
  // clang-format on

  bool success = pipeline.process(source, result);
  CHECK(success == true);
  CHECK(source == result);
}

TEST_CASE("case: vp8-mp4 convert", "[pipeline][vp8][mp4]") {
  auto source = initSource(0xAF);

  // clang-format off
  MP4Packet middle;
  auto pipeline = codecs::vp8enc() | muldex::vp8ToMP4();

  RawChunk result;
  auto pipeline2 = muldex::MP4Tovp8() | codecs::vp8dec();
  // clang-format on

  bool checkpoint1 = pipeline.process(source, middle);
  CHECK(checkpoint1 == true);

  bool checkpoint2 = pipeline2.process(middle, result);
  CHECK(checkpoint2 == true);

  CHECK(source == result);
}
#endif
