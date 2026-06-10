#ifndef _MEDIA_PIPE_H_
#define _MEDIA_PIPE_H_
#include <utility>

template <typename FROM, typename TO>
class Pipe {
 public:
  bool process(FROM in, TO& out) const {}
};

// clang-format off
template <typename SRC, typename MID, typename DST>
Pipe<SRC, DST> operator|(Pipe<SRC, MID> src, Pipe<MID, DST> sink) {
}
// clang-format on

#endif