#ifndef _MEDIA_PIPE_H_
#define _MEDIA_PIPE_H_
#include <memory>
#include <utility>


template <typename FROM, typename TO>
class Pipe {
 private:
  struct Pipeable {
    virtual ~Pipeable() = default;
    virtual bool process(FROM in, TO& out) const = 0;
  };

  template <typename Functor>
  struct Model : Pipeable {
    Functor fn;

    explicit Model(Functor f) : fn(std::move(f)) {}

    bool process(FROM in, TO& out) const override { return fn(in, out); }
  };
  std::shared_ptr<const Pipeable> pimpl_;

 public:
 
  template <typename F>
  Pipe(F f)
      : pimpl_(std::make_shared<Model<F>>(std::forward<F>(f))) {}

  bool process(FROM in, TO& out) const {
    if (!pimpl_) return false;
    return pimpl_->process(std::move(in), out);
  }
};

// clang-format off
template <typename SRC, typename MID, typename DST>
Pipe<SRC, DST> operator|(Pipe<SRC, MID> src, Pipe<MID, DST> sink) {
  return Pipe<SRC, DST>(
    [src = std::move(src), sink = std::move(sink)](SRC in, DST& out) -> bool {
    MID mid;
    return src.process(std::move(in), mid) && sink.process(std::move(mid), out);
  });
}
// clang-format on

#endif