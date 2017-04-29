#ifndef SCHEDULER_PROBLEM_H_
#define SCHEDULER_PROBLEM_H_

#include <iostream>

namespace scheduler {

typedef int group_t;
typedef int prof_t;
typedef int fatigue_t;
typedef int day_t;
typedef int class_time_t;

const group_t kMaxGroup = 60;
const prof_t kMaxProf = 60;
const int kDaysPerWeek = 6;
const int kClassesPerDay = 7;

struct Problem {
  int num_rooms;
  int num_groups;
  int num_profs;
  int num_classes[kMaxGroup + 1][kMaxProf + 1];

  friend std::istream& operator>>(std::istream& in, Problem& problem);
};

}  // namespace scheduler

#endif  // SCHEDULER_PROBLEM_H_
