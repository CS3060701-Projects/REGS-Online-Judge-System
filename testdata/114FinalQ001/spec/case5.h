#ifndef _CASE_H_
#define _CASE_H_
#define SECRET "@SECRET@"
#include "media_type.h"
#include "test.h"

TEST_CASE("case: h264-mp4 convert", "[pipeline][h264][mp4]") {
  auto source = initSource(0x01);
  RawChunk result;

  // clang-format off
  auto pipeline = codecs::h264enc() 
    | muldex::h264ToMP4()
    | muldex::MP4Toh264()
    | codecs::h264dec();
  // clang-format on

  bool success = pipeline.process(source, result);
  CHECK(success == true);
  CHECK(source.data == result.data);
}

TEST_CASE("case: vp8-mp4 convert", "[pipeline][vp8][mp4]") {
  auto source = initSource(0x26);

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

  CHECK(source.data == result.data);
}

TEST_CASE("case: vp9-MKV convert", "[pipeline][vp9][mkv]") {
  auto source = initSource(0x10);

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

  CHECK(source.data == result.data);
}

TEST_CASE("case: h265-MKV convert", "[pipeline][h265][mkv]") {
  auto source = initSource(0x18);
  RawChunk result;

  // clang-format off
  auto pipeline = codecs::h265enc() 
    | muldex::h265ToMKV()
    | muldex::MKVToh265()
    | codecs::h265dec();
  // clang-format on

  bool success = pipeline.process(source, result);
  CHECK(success == true);
  CHECK(source.data == result.data);
}

TEST_CASE("case: h266-flv convert", "[pipeline][av1][flv]") {
  auto source = initSource(0x87);

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

  CHECK(source.data == result.data);
}

TEST_CASE("case: av1-flv convert", "[pipeline][av1][flv]") {
  auto source = initSource(0xAF);
  RawChunk result;

  // clang-format off
  auto pipeline = codecs::av1enc()
    | muldex::av1ToFLV()
    | muldex::FLVToav1()
    | codecs::av1dec();
  // clang-format on

  bool success = pipeline.process(source, result);
  CHECK(success == true);
  CHECK(source.data == result.data);
}

namespace filter {
inline Pipe<H264Frame, H264Frame> myfilter() {
  return Pipe<H264Frame, H264Frame>([](H264Frame in, H264Frame& out) {
    out.data[3] = 0x00;
    out.data[9] = 0x15;
    return true;
  });
}
}  // namespace filter

TEST_CASE("case: custom filter corruption", "[pipeline][filter]") {
  auto source = initSource(0xCC);
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