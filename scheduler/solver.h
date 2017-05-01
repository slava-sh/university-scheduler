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
  struct DayStats {
    fatigue_t fatigue;
    bool has_skips;
    day_t min_class;
    day_t max_class;
  };

  struct State : public Solution {
    // Inherit instead of embedding to save typing.
    group_t prof_schedule[kMaxProf + 1][kDaysPerWeek + 1][kClassesPerDay + 1];
    int num_free_rooms[kDaysPerWeek + 1][kClassesPerDay + 1];
    DayStats group_stats[kMaxGroup + 1][kDaysPerWeek + 1];
    DayStats prof_stats[kMaxProf + 1][kDaysPerWeek + 1];
  };

  const int kMaxIdleSteps = 1000000;

  State SolveNaive(const std::shared_ptr<Problem> &problem);
  fatigue_t PartialFatigue(const int *schedule, const Solver::DayStats &stats);
  fatigue_t GroupFatigue(const State &state, group_t group, day_t day);
  fatigue_t ProfFatigue(const State &state, prof_t prof, day_t day);
};

}  // namespace scheduler

#endif  // SCHEDULER_SOLVER_H_
