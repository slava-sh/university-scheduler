#ifndef SCHEDULER_TIME_LIMITED_SOLVER_H_
#define SCHEDULER_TIME_LIMITED_SOLVER_H_

#include <memory>
#include "scheduler/solver.h"

namespace scheduler {

class TimeLimitedSolver : public Solver {
 public:
  explicit TimeLimitedSolver(std::chrono::nanoseconds time_limit)
      : time_limit_(time_limit) {}

  Solution Solve(const std::shared_ptr<Problem>& problem) override;
  bool ShouldStop() override;

 private:
  typedef std::chrono::steady_clock clock;
  std::chrono::nanoseconds time_limit_;
  clock::time_point start_;
};

}  // namespace scheduler

#endif  // SCHEDULER_TIME_LIMITED_SOLVER_H_
