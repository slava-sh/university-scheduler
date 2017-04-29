#include "scheduler/time_limited_solver.h"

namespace scheduler {

Solution TimeLimitedSolver::Solve(const std::shared_ptr<Problem>& problem) {
  start_ = clock::now();
  return Solver::Solve(problem);
}

bool TimeLimitedSolver::ShouldStop() {
  auto elapsed = clock::now() - start_;
  return elapsed >= time_limit_;
}

}  // namespace scheduler
