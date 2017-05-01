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
}

void Solver::State::Day::RemoveClass(class_time_t time) { schedule_[time] = 0; }

fatigue_t Solver::State::Day::Fatigue() {
  class_time_t max_time = 0;
  for (class_time_t time = kClassesPerDay; time > 0; --time) {
    if (schedule_[time] != 0) {
      max_time = time;
      break;
    }
  }
  if (max_time == 0) {
    return 0;
  }
  class_time_t min_time = 0;
  for (class_time_t time = 1; time <= kClassesPerDay; ++time) {
    if (schedule_[time] != 0) {
      min_time = time;
      break;
    }
  }
  return square(2 + max_time - min_time + 1);
}

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
      auto d1 = 1 + Random(kDaysPerWeek);
      auto c1 = 1 + Random(kClassesPerDay);
      auto p = 1 + Random(problem->num_profs);
      auto g = state.prof[p][d1].GetClass(c1);
      if (g == 0) {
        continue;
      }
      auto d2 = 1 + Random(kDaysPerWeek);
      if (d2 == d1) {
        continue;
      }
      auto c2 = 1 + Random(kClassesPerDay);
      if (state.num_free_rooms[d2][c2] == 0 || state.prof[p][d2].HasClass(c2) ||
          state.group[g][d2].HasClass(c2)) {
        continue;
      }

      if (1 < c1 && c1 < kClassesPerDay) {
        auto group_will_have_empty_slot = state.group[g][d1].HasClass(c1 - 1) &&
                                          state.group[g][d1].HasClass(c1 + 1);
        if (group_will_have_empty_slot) {
          continue;
        }
        auto prof_will_have_empty_slot = state.prof[p][d1].HasClass(c1 - 1) &&
                                         state.prof[p][d1].HasClass(c1 + 1);
        if (prof_will_have_empty_slot) {
          continue;
        }
      }

      // Apply swap.
      auto new_group1 = state.group[g][d1];
      auto new_group2 = state.group[g][d2];
      auto new_prof1 = state.prof[p][d1];
      auto new_prof2 = state.prof[p][d2];
      new_group1.RemoveClass(c1);
      new_group2.AddClass(c2, p);
      new_prof1.RemoveClass(c1);
      new_prof2.AddClass(c2, g);

      auto new_fatigue = state.fatigue;
      new_fatigue -= state.group[g][d1].Fatigue();
      new_fatigue -= state.prof[p][d1].Fatigue();
      new_fatigue -= state.group[g][d2].Fatigue();
      new_fatigue -= state.prof[p][d2].Fatigue();
      new_fatigue += new_group1.Fatigue();
      new_fatigue += new_prof1.Fatigue();
      new_fatigue += new_group2.Fatigue();
      new_fatigue += new_prof2.Fatigue();

      if (new_fatigue <= state.fatigue) {
        // Accept swap.
        state.group[g][d1] = new_group1;
        state.group[g][d2] = new_group2;
        state.prof[p][d1] = new_prof1;
        state.prof[p][d2] = new_prof2;
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
