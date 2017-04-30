#include "scheduler/problem.h"

namespace scheduler {

std::istream &operator>>(std::istream &in, Problem &problem) {
  in >> problem.num_groups;
  in >> problem.num_profs;
  in >> problem.num_rooms;
  for (group_t group = 1; group <= problem.num_groups; ++group) {
    for (prof_t prof = 1; prof <= problem.num_profs; ++prof) {
      in >> problem.num_classes[group][prof];
    }
  }
  return in;
}

}  // namespace scheduler
