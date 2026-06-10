#ifndef _MEDIA_H_
#define _MEDIA_H_
#include <array>
#include <cstdint>
#include <iomanip>
#include <iostream>

#include "media_pipe.h"
#include "test.h"

using Buffer = std::array<uint8_t, 10>;

// Raw input
struct RawChunk {
  Buffer data;

  bool operator==(const RawChunk& other) const { return data == other.data; }
};

namespace Catch {
template <>
struct StringMaker<RawChunk> {
  static std::string convert(RawChunk const& chunk) {
    std::ostringstream oss;
    oss << "\"";
    for (uint8_t byte : chunk.data) {
      oss << "0x" << std::hex << std::setw(2) << std::setfill('0')
          << std::uppercase << static_cast<int>(byte);
    }
    oss << "\"";
    return oss.str();
  }
};
}  // namespace Catch

// clang-format off
// Different encoder
struct H264Frame: public RawChunk {};
struct H265Frame: public RawChunk {};
struct H266Frame: public RawChunk {};
struct VP8Frame: public RawChunk {};
struct VP9Frame: public RawChunk {};
struct AV1Frame: public RawChunk {};

// Different container
struct MP4Packet: public RawChunk {};
struct MKVPacket: public RawChunk {};
struct FLVPacket: public RawChunk {};

// clang-format on

// Utility

inline RawChunk initSource(uint8_t pattern) {
  RawChunk src;
  src.data.fill(pattern);
  return src;
}

inline void convert(const Buffer& in, Buffer& out, uint8_t key) {
  for (std::size_t i = 0; i < in.size(); ++i) {
    out[i] = key ^ in[i];
  }
}

#define PROCESS(FROM, TO, KEY)                   \
  Pipe<FROM, TO> {                               \
    return Pipe<FROM, TO>([](FROM in, TO& out) { \
      convert(in.data, out.data, KEY);           \
      return true;                               \
    });                                          \
  }

namespace codecs {
inline auto h264enc() -> PROCESS(RawChunk, H264Frame, 0x65);
inline auto h264dec() -> PROCESS(H264Frame, RawChunk, 0x65);

inline auto h265enc() -> PROCESS(RawChunk, H265Frame, 0x66);
inline auto h265dec() -> PROCESS(H265Frame, RawChunk, 0x66);

inline auto h266dec() -> PROCESS(H266Frame, RawChunk, 0x67);
inline auto h266enc() -> PROCESS(RawChunk, H266Frame, 0x67);

inline auto vp8enc() -> PROCESS(RawChunk, VP8Frame, 0x95);
inline auto vp8dec() -> PROCESS(VP8Frame, RawChunk, 0x95);

inline auto vp9enc() -> PROCESS(RawChunk, VP9Frame, 0x96);
inline auto vp9dec() -> PROCESS(VP9Frame, RawChunk, 0x96);

inline auto av1enc() -> PROCESS(RawChunk, AV1Frame, 0x97);
inline auto av1dec() -> PROCESS(AV1Frame, RawChunk, 0x97);
}  // namespace codecs

namespace muldex {
// MP4
inline auto h264ToMP4() -> PROCESS(H264Frame, MP4Packet, 0x11);
inline auto MP4Toh264() -> PROCESS(MP4Packet, H264Frame, 0x11);

inline auto vp8ToMP4() -> PROCESS(VP8Frame, MP4Packet, 0x12);
inline auto MP4Tovp8() -> PROCESS(MP4Packet, VP8Frame, 0x12);

// MKV
inline auto h265ToMKV() -> PROCESS(H265Frame, MKVPacket, 0x21);
inline auto MKVToh265() -> PROCESS(MKVPacket, H265Frame, 0x21);

inline auto vp9ToMKV() -> PROCESS(VP9Frame, MKVPacket, 0x22);
inline auto MKVTovp9() -> PROCESS(MKVPacket, VP9Frame, 0x22);

// FLV
inline auto h266ToFLV() -> PROCESS(H266Frame, FLVPacket, 0x31);
inline auto FLVToh266() -> PROCESS(FLVPacket, H266Frame, 0x31);

inline auto av1ToFLV() -> PROCESS(AV1Frame, FLVPacket, 0x32);
inline auto FLVToav1() -> PROCESS(FLVPacket, AV1Frame, 0x32);
}  // namespace muldex

#endif