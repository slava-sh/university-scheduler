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
  struct State : public Solution {
    // Inherit instead of embedding to save typing.
    group_t prof_schedule[kMaxProf + 1][kDaysPerWeek + 1][kClassesPerDay + 1];
    int num_free_rooms[kDaysPerWeek + 1][kClassesPerDay + 1];
    fatigue_t group_fatigue[kMaxGroup + 1][kDaysPerWeek + 1];
    fatigue_t prof_fatigue[kMaxProf + 1][kDaysPerWeek + 1];
  };

  State SolveNaive(const std::shared_ptr<Problem>& problem);
  fatigue_t PartialFatigue(const int *schedule);
  fatigue_t GroupFatigue(const State &state, group_t group, day_t day);
  fatigue_t ProfFatigue(const State &state, prof_t prof, day_t day);
};

}  // namespace scheduler

#endif  // SCHEDULER_SOLVER_H_
