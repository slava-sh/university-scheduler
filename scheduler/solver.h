#ifndef SCHEDULER_SOLVER_H_
#define SCHEDULER_SOLVER_H_

#include <memory>
#include "scheduler/solution.h"

namespace scheduler {

class Solver {
 public:
  virtual Solution Solve(const std::shared_ptr<Problem> &problem);

 protected:
  virtual bool ShouldStop() = 0;

 private:
  template <typename schedule_t>
  struct DayState {
    schedule_t schedule[kClassesPerDay + 1];
    fatigue_t fatigue;
  };

  struct State {
    fatigue_t fatigue;
    DayState<prof_t> group[kMaxGroup + 1][kDaysPerWeek + 1];
    DayState<group_t> prof[kMaxProf + 1][kDaysPerWeek + 1];
    int num_free_rooms[kDaysPerWeek + 1][kClassesPerDay + 1];
  };

  const int kMaxIdleSteps = 1000000;

  State SolveNaive(const std::shared_ptr<Problem> &problem);

  template <typename schedule_t>
  fatigue_t Fatigue(const DayState<schedule_t> &day_state);
};

}  // namespace scheduler

#endif  // SCHEDULER_SOLVER_H_
