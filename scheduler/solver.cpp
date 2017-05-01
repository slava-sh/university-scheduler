#include "scheduler/solver.h"
#include <unordered_map>
#include <utility>
#include <vector>

namespace scheduler {

int square(int x) { return x * x; }

void Solver::UpdateStats(const int *schedule, Solver::DayStats &stats) {
  stats.min_class = 0;
  for (class_time_t time = 1; time <= kClassesPerDay; ++time) {
    if (schedule[time] != 0) {
      stats.min_class = time;
      break;
    }
  }
  if (stats.min_class == 0) {
    stats.max_class = 0;
    stats.has_skips = false;
  } else {
    for (class_time_t time = kClassesPerDay; time > 0; --time) {
      if (schedule[time] != 0) {
        stats.max_class = time;
        break;
      }
    }
    stats.has_skips = false;
    for (class_time_t time = stats.min_class + 1; time < stats.max_class;
         ++time) {
      if (schedule[time] == 0) {
        stats.has_skips = true;
        break;
      }
    }
  }
  stats.fatigue = PartialFatigue(stats);
}

fatigue_t Solver::PartialFatigue(const Solver::DayStats &stats) {
  if (stats.min_class == 0) {
    return 0;
  }
  return square(2 + stats.max_class - stats.max_class + 1);
}

fatigue_t Solver::GroupFatigue(const State &state, group_t group, day_t day) {
  return PartialFatigue(state.group_stats[group][day]);
}

fatigue_t Solver::ProfFatigue(const State &state, prof_t prof, day_t day) {
  return PartialFatigue(state.prof_stats[prof][day]);
}

namespace {

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

}  // namespace

Solver::State Solver::SolveNaive(const std::shared_ptr<Problem> &problem) {
  State state = {};
  state.problem = problem;
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
        state.group_schedule[group][day][time] = prof;
        state.prof_schedule[prof][day][time] = group;
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
      UpdateStats(state.group_schedule[group][day],
                  state.group_stats[group][day]);
      state.fatigue += state.group_stats[group][day].fatigue;
    }
    for (prof_t prof = 1; prof <= problem->num_profs; ++prof) {
      UpdateStats(state.prof_schedule[prof][day], state.prof_stats[prof][day]);
      state.fatigue += state.prof_stats[prof][day].fatigue;
    }
  }
  return state;
}

int Random(int n) {
  // Note: This generator is not uniform but is probably faster.
  return std::rand() % n;
}

bool RandomBool() { return Random(2) == 0; }

Solution Solver::Solve(const std::shared_ptr<Problem> &problem) {
  auto state = SolveNaive(problem);
  Solution best_solution(state);
  int idle_steps = 0;
  while (!ShouldStop()) {
    if (idle_steps == kMaxIdleSteps) {
      state = SolveNaive(problem);
      idle_steps = 0;
      continue;
    }
    for (int t = 0; t < 50; ++t) {
      // Generate a swap.
      // TODO: if (state.has_skips) {
      auto d1 = 1 + Random(kDaysPerWeek);
      auto c1 = 1 + Random(kClassesPerDay);
      auto p = 1 + Random(problem->num_profs);
      auto g = state.prof_schedule[p][d1][c1];
      if (g == 0) {
        continue;
      }
      auto d2 = 1 + Random(kDaysPerWeek);
      auto c2 = 1 + Random(kClassesPerDay);
      if (state.num_free_rooms[d2][c2] == 0 ||
          state.prof_schedule[p][d2][c2] != 0 ||
          state.group_schedule[g][d2][c2] != 0) {
        continue;
      }

      if (1 < c1 && c1 < kClassesPerDay) {
        auto group_will_have_empty_slot =
            state.group_schedule[g][d1][c1 - 1] != 0 &&
            state.group_schedule[g][d1][c1 + 1] != 0;
        if (group_will_have_empty_slot) {
          continue;
        }
        auto prof_will_have_empty_slot =
            state.prof_schedule[p][d1][c1 - 1] != 0 &&
            state.prof_schedule[p][d1][c1 + 1] != 0;
        if (prof_will_have_empty_slot) {
          continue;
        }
      }
      /*
      } else {
        auto g = 1 + Random(problem->num_groups);

        auto d2 = 1 + Random(kDaysPerWeek);
        auto c2_is_min = RandomBool();
        auto c2 = c2_is_min ? state.group_stats[g][d2].min_class - 1
                            : state.group_stats[g][d2].max_class + 1;
        if (!(1 <= c2 && c2 <= kClassesPerDay)) {
          continue;
        }

        auto d1 = 1 + Random(kDaysPerWeek);
        auto c1_is_min = RandomBool();
        auto c1 = c1_is_min ? state.group_stats[g][d1].min_class
                            : state.group_stats[g][d1].max_class;
        auto p = state.group_schedule[g][d1][c1];

        auto prof_will_have_empty_slot =
            (1 < c1 && c1 < kClassesPerDay &&
             state.prof_schedule[p][d1][c1 - 1] != 0 &&
             state.prof_schedule[p][d1][c1 + 1] != 0) ||
            !(state.prof_stats[p][d2].min_class == 0 ||
              c2 == state.prof_stats[p][d2].min_class - 1 ||
              c2 == state.prof_stats[p][d2].max_class + 1);
        if (prof_will_have_empty_slot) {
          continue;
        }
      }
       */

      auto prev_fatigue = state.fatigue;
      auto prev_group_stats1 = state.group_stats[g][d1];
      auto prev_group_stats2 = state.group_stats[g][d2];
      auto prev_prof_stats1 = state.prof_stats[p][d1];
      auto prev_prof_stats2 = state.prof_stats[p][d2];

      // Apply swap.
      state.fatigue -= state.group_stats[g][d1].fatigue;
      state.fatigue -= state.prof_stats[p][d1].fatigue;
      if (d2 != d1) {
        state.fatigue -= state.group_stats[g][d2].fatigue;
        state.fatigue -= state.prof_stats[p][d2].fatigue;
      }
      state.group_schedule[g][d1][c1] = 0;
      state.group_schedule[g][d2][c2] = p;
      state.prof_schedule[p][d1][c1] = 0;
      state.prof_schedule[p][d2][c2] = g;
      RemoveClass(state.group_stats[g][d1], c1, state.group_schedule[g][d1]);
      RemoveClass(state.prof_stats[p][d1], c1, state.prof_schedule[p][d1]);
      AddClass(state.group_stats[g][d2], c2);
      AddClass(state.prof_stats[p][d2], c2);
      state.group_stats[g][d1].fatigue = GroupFatigue(state, g, d1);
      state.prof_stats[p][d1].fatigue = ProfFatigue(state, p, d1);
      state.fatigue += state.group_stats[g][d1].fatigue;
      state.fatigue += state.prof_stats[p][d1].fatigue;
      if (d2 != d1) {
        state.group_stats[g][d2].fatigue = GroupFatigue(state, g, d2);
        state.prof_stats[p][d2].fatigue = ProfFatigue(state, p, d2);
        state.fatigue += state.group_stats[g][d2].fatigue;
        state.fatigue += state.prof_stats[p][d2].fatigue;
      }

      if (state.fatigue <= prev_fatigue) {
        // Accept swap.
        ++state.num_free_rooms[d1][c1];
        --state.num_free_rooms[d2][c2];
        if (state.fatigue < best_solution.fatigue) {
          best_solution = Solution(state);
        }
      } else {
        // Reject swap.
        state.group_schedule[g][d1][c1] = p;
        state.group_schedule[g][d2][c2] = 0;
        state.prof_schedule[p][d1][c1] = g;
        state.prof_schedule[p][d2][c2] = 0;
        state.fatigue = prev_fatigue;
        state.group_stats[g][d1] = prev_group_stats1;
        state.prof_stats[p][d1] = prev_prof_stats1;
        if (d2 != d1) {
          state.group_stats[g][d2] = prev_group_stats2;
          state.prof_stats[p][d2] = prev_prof_stats2;
        }
      }

      if (state.fatigue == prev_fatigue) {
        ++idle_steps;
      } else {
        idle_steps = 0;
      }

      break;
    }
  }
  return best_solution;
}

void Solver::AddClass(Solver::DayStats &stats, class_time_t time) {
  if (stats.min_class == 0) {
    stats.min_class = time;
    stats.max_class = time;
    return;
  }
  if (time < stats.min_class) {
    stats.min_class = time;
  } else if (time > stats.max_class) {
    stats.max_class = time;
  }
  // TODO: Maybe update `has_skips`.
}

void Solver::RemoveClass(Solver::DayStats &stats, class_time_t time,
                         const int *schedule) {
  if (stats.min_class == stats.max_class) {
    stats.min_class = 0;
    stats.max_class = 0;
    return;
  }
  if (time == stats.min_class) {
    do {
      ++stats.min_class;
    } while (schedule[stats.min_class] == 0);
  } else if (time == stats.max_class) {
    do {
      --stats.max_class;
    } while (schedule[stats.max_class] == 0);
  }
  // TODO: Maybe update `has_skips`.
}

}  // namespace scheduler
