#include "scheduler/solution.h"

namespace scheduler {

std::ostream &operator<<(std::ostream &out, const Solution &solution) {
  out << solution.fatigue << "\n";
  auto problem = solution.problem;
  for (group_t group = 1; group <= problem->num_groups; group++) {
    out << "\n";
    for (class_time_t time = 1; time <= kClassesPerDay; time++) {
      for (day_t day = 1; day <= kDaysPerWeek; day++) {
        if (day != 1) {
          out << " ";
        }
        out << solution.group_schedule[group][day][time];
      }
      out << "\n";
    }
  }
  return out;
}

}  // namespace scheduler
