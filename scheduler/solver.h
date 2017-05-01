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
  struct State {
    class Day {
     public:
      int GetClass(class_time_t time) const;
      bool HasClass(class_time_t time) const;
      void AddClass(class_time_t time, int value);
      void RemoveClass(class_time_t time);
      fatigue_t Fatigue();
      bool HasClasses() const;
      int GetMinClass() const;
      int GetMaxClass() const;
      bool HasSkips() const;

     private:
      int schedule_[kClassesPerDay + 1] = {};
      class_time_t min_class_ = 0;
      class_time_t max_class_ = 0;
      int num_classes_ = 0;
    };

    fatigue_t fatigue = 0;
    Day group[kMaxGroup + 1][kDaysPerWeek + 1];
    Day prof[kMaxProf + 1][kDaysPerWeek + 1];
    int num_free_rooms[kDaysPerWeek + 1][kClassesPerDay + 1] = {};
  };

  const int kMaxIdleSteps = 1000000;

  State SolveNaive(const std::shared_ptr<Problem> &problem);
};

}  // namespace scheduler

#endif  // SCHEDULER_SOLVER_H_
