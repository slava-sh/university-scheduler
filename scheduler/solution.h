#ifndef SCHEDULER_SOLUTION_H_
#define SCHEDULER_SOLUTION_H_

#include <iostream>
#include <memory>
#include "scheduler/problem.h"

namespace scheduler {

struct Solution {
  std::shared_ptr<Problem> problem;
  fatigue_t fatigue;
  prof_t group_schedule[kMaxGroup + 1][kDaysPerWeek + 1][kClassesPerDay + 1];

  friend std::ostream &operator<<(std::ostream &out, const Solution &solution);
};

}  // namespace scheduler

#endif  // SCHEDULER_SOLUTION_H_
