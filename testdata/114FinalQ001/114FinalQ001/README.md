---
author: Chen, YanHau
date: 2026/05/17
category: Tempalte, Operator Overloading
difficulty: Hard
expect-time: 60-90min
---
# Stream Pipe

<div style="display: flex; flex-direction: row; justify-content: space-between">
  <div>Difficulty: Hard</div>
  <div>Expect Time: 60-90min</div>
  <div>Author: Chen, YanHau</div>
</div>

Streaming services are a very common application in modern society. Referring to the architecture diagrams of systems such as GStreamer or FFmpeg, they typically involve a series of operations such as demuxing, decoding, encoding, and muxing, which package raw data into media containers.

<img style="display: block; margin: 0 auto" src="https://upload.wikimedia.org/wikipedia/commons/a/a5/GStreamer_example_pipeline.svg" width="75%" alt="Gstreamer" />

In this problem, we will implement a pipeline that simulates a media processing workflow. The following are several core components:

- Buffer: The type used to store the actual data, defined as an alias of `std::array<uint8_t, 10>`.
- RawChunk: The basic input type defined by the problem setter; all transformations begin from this type.
- `*Frame` and `*Packet`: Struct includes `{ Buffer data }` member.

## Codecs & Muldex

In `media.h`, six major types are defined: `H264Frame`, `H265Frame`, `H266Frame`, `VP8Frame`, `VP9Frame`, and `AV1Frame`, corresponding to the `codecs::*enc` and `codecs::*dec` series of functions. Similarly, muxers and demuxers implement a series of `<Frame>To<Packet>` or `<Packet>To<Frame>` functions. These methods that attempt to convert A into B are referred to as *pipeable functions*. A *pipeable function* returns an instance of the class `Pipe<FROM, TO>`.

```cpp
// For example, `h264enc` returns an instance of `Pipe<RawChunk, H264Frame>`.
// h264enc is a `pipeable function`
inline Pipe<RawChunk, H264Frame> h264enc() {
  return Pipe<RawChunk, H264Frame>([](RawChunk in, H264Frame& out) {
    process(in.data, out.data, 0x65);
    return true;
  });
}
```

## Pipe

For class `Pipe<FROM, TO>`, you need to implement two key methods:

- `Pipe::pipe`, which allows registering a function of the form bool(FROM in, TO& out).
- `bool Pipe::process(FROM in, TO& out) const`, which attempts to invoke the pipeable function defined in `media_type.h`. It returns `true` if the conversion succeeds; otherwise, it returns `false`.

You also need to overloading `operator|(PipeA, PipeB)` . It takes two Pipe instances, and if they can be connected, it attempts to pass data through the combined pipeline for conversion.

```cpp
template <typename SRC, typename MID, typename DST>
auto operator|(Pipe<SRC, MID> src, Pipe<MID, DST> sink) -> Pipe<SRC, DST> { }
```

In case*.h, a set of pipelines is define:

```cpp
auto pipeline = codecs::h264enc()
  | muldex::h264ToMP4()
  | muldex::MP4Toh264();

RawChunk chunk;
H264Frame frame;
pipeline.process(chunk, frame);
```

The data is expected to flow through the pipeline:

`Pipe<RawChunk, H264Frame> | Pipe<H264Frame, MP4Packet> | Pipe<MP4Packet, H264Frame>`.

Since `operator|` is used to chain multiple pipes, the final result should be `Pipe<RawChunk, H264Frame>`. Therefore, pipeline.process takes RawChunk as input and produces H264Frame as output.

Complete the missing implementation in `media_pipe.h` so that all test cases can run successfully.

<div style="page-break-after: always;"></div>

## Input

1. Complete the required class implementation. The following files already contain some implementations:

    - `media.h`
    - `media_pipe.h`

2. The Online Judge system will replace the following files. **DON'T** write implemention in those file:

    - `entrypoint.cpp`
    - `test.h`
    - `case.h`

3. The test cases for student available in `case*.h`.

<div style="page-break-after: always;"></div>

## Output

Online Judge will run the test:

1. **DON' T** output any information to the terminal.
2. When any test failed, Online Judge will try to list error messages.
