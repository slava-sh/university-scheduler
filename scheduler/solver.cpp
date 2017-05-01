#include "scheduler/solver.h"
#include <unordered_map>
#include <utility>
#include <vector>

namespace scheduler {

int square(int x) { return x * x; }

int Solver::State::Day::GetClass(class_time_t time) const {
  return schedule_[time];
}

bool Solver::State::Day::HasClass(class_time_t time) const {
  return schedule_[time] != 0;
}

void Solver::State::Day::AddClass(class_time_t time, int value) {
  schedule_[time] = value;
  if (num_classes_ == 0) {
    min_class_ = time;
    max_class_ = time;
  } else if (time < min_class_) {
    min_class_ = time;
  } else if (time > max_class_) {
    max_class_ = time;
  }
  ++num_classes_;
}

void Solver::State::Day::RemoveClass(class_time_t time) {
  schedule_[time] = 0;
  --num_classes_;
  if (num_classes_ == 0) {
    min_class_ = 0;
    max_class_ = 0;
  } else if (time == min_class_) {
    do {
      ++min_class_;
    } while (schedule_[min_class_] == 0);
  } else if (time == max_class_) {
    do {
      --max_class_;
    } while (schedule_[max_class_] == 0);
  }
}

fatigue_t Solver::State::Day::Fatigue() {
  if (num_classes_ == 0) {
    return 0;
  }
  return square(2 + max_class_ - min_class_ + 1);
}

bool Solver::State::Day::HasSkips() const {
  return HasClasses() && max_class_ - min_class_ + 1 == num_classes_;
}

int Solver::State::Day::GetMinClass() const { return min_class_; }

int Solver::State::Day::GetMaxClass() const { return max_class_; }

bool Solver::State::Day::HasClasses() const { return num_classes_ != 0; }

struct hash {
  explicit hash(size_t seed) : seed_(seed) {}

  template <class T1, class T2>
  std::size_t operator()(const std::pair<T1, T2> &p) const {
    auto h1 = std::hash<T1>()(p.first);
    auto h2 = std::hash<T2>()(p.second);
    return h1 ^ h2 ^ seed_;
  }

 private:
  size_t seed_;
};

Solver::State Solver::SolveNaive(const std::shared_ptr<Problem> &problem) {
  State state;
  auto seed = std::rand();
  std::unordered_map<std::pair<group_t, prof_t>, int, hash> classes_to_schedule(
      0, hash(seed));
  for (group_t group = 1; group <= problem->num_groups; ++group) {
    for (prof_t prof = 1; prof <= problem->num_profs; ++prof) {
      if (problem->num_classes[group][prof] == 0) {
        continue;
      }
      auto group_and_prof = std::make_pair(group, prof);
      classes_to_schedule[group_and_prof] = problem->num_classes[group][prof];
    }
  }
  for (day_t day = 1; day <= kDaysPerWeek; ++day) {
    for (class_time_t time = 1; time <= kClassesPerDay; ++time) {
      state.num_free_rooms[day][time] = problem->num_rooms;
      std::vector<bool> group_is_busy(problem->num_groups + 1);
      std::vector<bool> prof_is_busy(problem->num_profs + 1);
      for (auto it = classes_to_schedule.begin();
           it != classes_to_schedule.end();) {
        auto &group_and_prof = it->first;
        auto &group = group_and_prof.first;
        auto &prof = group_and_prof.second;
        if (group_is_busy[group] || prof_is_busy[prof]) {
          ++it;
          continue;
        }
        if (state.num_free_rooms[day][time] == 0) {
          break;
        }
        state.group[group][day].AddClass(time, prof);
        state.prof[prof][day].AddClass(time, group);
        group_is_busy[group] = true;
        prof_is_busy[prof] = true;
        --state.num_free_rooms[day][time];
        --it->second;
        if (it->second == 0) {
          it = classes_to_schedule.erase(it);
        } else {
          ++it;
        }
      }
    }
  }

  state.fatigue = 0;
  for (day_t day = 1; day <= kDaysPerWeek; ++day) {
    for (group_t group = 1; group <= problem->num_groups; ++group) {
      state.fatigue += state.group[group][day].Fatigue();
    }
    for (prof_t prof = 1; prof <= problem->num_profs; ++prof) {
      state.fatigue += state.prof[prof][day].Fatigue();
    }
  }
  return state;
}

int Random(int n) {
  // Note: This generator is not uniform but is probably faster.
  return std::rand() % n;
}

int RandomBetween(int a, int b) { return a + Random(b - a + 1); }

bool RandomBool() { return Random(2) == 0; }

Solution Solver::Solve(const std::shared_ptr<Problem> &problem) {
  auto state = SolveNaive(problem);
  int idle_steps = 0;
  while (!ShouldStop()) {
    /*
    if (idle_steps == kMaxIdleSteps) {
      state = SolveNaive(problem);
      idle_steps = 0;
      continue;
    }
     */
    for (int t = 0; t < 50; ++t) {
      // Generate a swap.
      auto g = 1 + Random(problem->num_groups);
      auto d1 = 1 + Random(kDaysPerWeek);
      auto &group1 = state.group[g][d1];
      if (!group1.HasClasses()) {
        continue;
      }
      auto c1 = group1.HasSkips() ? 1 + Random(kClassesPerDay)
                                  : (RandomBool() ? group1.GetMinClass()
                                                  : group1.GetMaxClass());
      auto p = group1.GetClass(c1);
      if (p == 0) {
        continue;
      }
      auto &prof1 = state.prof[p][d1];
      if ((1 < c1 && c1 < kClassesPerDay) &&
          (prof1.HasClass(c1 - 1) && prof1.HasClass(c1 + 1))) {
        continue;
      }
      auto d2 = 1 + Random(kDaysPerWeek);
      if (d2 == d1) {
        continue;
      }
      auto &group2 = state.group[g][d2];
      auto &prof2 = state.prof[p][d2];
      auto c2 = group2.HasSkips() ? 1 + Random(kClassesPerDay)
                                  : (RandomBool() ? group2.GetMinClass() - 1
                                                  : group2.GetMaxClass() + 1);
      if (!(1 <= c2 && c2 <= kClassesPerDay) || group2.HasClass(c2) ||
          prof2.HasClass(c2) || state.num_free_rooms[d2][c2] == 0) {
        continue;
      }

      // Apply swap.
      auto new_group1 = group1;
      auto new_group2 = group2;
      auto new_prof1 = prof1;
      auto new_prof2 = prof2;
      new_group1.RemoveClass(c1);
      new_prof1.RemoveClass(c1);
      new_group2.AddClass(c2, p);
      new_prof2.AddClass(c2, g);

      auto new_fatigue = state.fatigue;
      new_fatigue -= group1.Fatigue();
      new_fatigue -= group2.Fatigue();
      new_fatigue -= prof1.Fatigue();
      new_fatigue -= prof2.Fatigue();
      new_fatigue += new_group1.Fatigue();
      new_fatigue += new_group2.Fatigue();
      new_fatigue += new_prof1.Fatigue();
      new_fatigue += new_prof2.Fatigue();

      if (new_fatigue <= state.fatigue) {
        // Accept swap.
        group1 = new_group1;
        group2 = new_group2;
        prof1 = new_prof1;
        prof2 = new_prof2;
        state.fatigue = new_fatigue;
        ++state.num_free_rooms[d1][c1];
        --state.num_free_rooms[d2][c2];
        /* Our state always represents the best solution known.
        if (state.fatigue < best_solution.fatigue) {
          best_solution = Solution(state);
        }
         */
      }

      if (new_fatigue >= state.fatigue) {
        ++idle_steps;
      } else {
        idle_steps = 0;
      }

      break;
    }
  }

  Solution solution;
  solution.problem = problem;
  solution.fatigue = state.fatigue;
  for (group_t group = 1; group <= problem->num_groups; ++group) {
    for (day_t day = 1; day <= kDaysPerWeek; ++day) {
      for (class_time_t time = 1; time <= kClassesPerDay; ++time) {
        solution.group_schedule[group][day][time] =
            state.group[group][day].GetClass(time);
      }
    }
  }
  return solution;
}

}  // namespace scheduler
