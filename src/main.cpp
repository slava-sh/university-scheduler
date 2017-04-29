#include <iostream>
#include <fstream>
#include <iomanip>
#include <memory>
#include <unordered_map>
#include <vector>
#include <chrono>

namespace {

const auto time_limit = std::chrono::seconds(10) - std::chrono::milliseconds(50);

typedef int group;
typedef int prof;
typedef int fatigue;
typedef int day;
typedef int class_time;

const group MaxGroup    = 60;
const prof MaxProf      = 60;
const int DaysPerWeek   = 6;
const int ClassesPerDay = 7;

struct Problem {
    int NumRooms;
    int NumGroups;
    int NumProfs;
    int NumClasses[MaxGroup + 1][MaxProf + 1];
};

struct Solution {
    std::shared_ptr<Problem> problem;
    fatigue Fatigue;
    prof GroupSchedule[MaxGroup + 1][DaysPerWeek + 1][ClassesPerDay + 1];

    Solution(std::shared_ptr<Problem> problem);
    Solution(const Solution& other) = default;
};

Solution::Solution(const std::shared_ptr<Problem> problem):
    problem(problem), GroupSchedule() {
}

struct State : public Solution {
    group ProfSchedule[MaxProf + 1][DaysPerWeek + 1][ClassesPerDay + 1];
    int NumFreeRooms[DaysPerWeek + 1][ClassesPerDay + 1];
    fatigue GroupFatigue[MaxGroup + 1][DaysPerWeek + 1];
    fatigue ProfFatigue[MaxProf + 1][DaysPerWeek + 1];

    State(std::shared_ptr<Problem> problem);

    fatigue partialFatigue(int* schedule);
    fatigue groupFatigue(group group, day day);
    fatigue profFatigue(prof prof, day day);
};

State::State(std::shared_ptr<Problem> problem):
    Solution(problem), ProfSchedule(), NumFreeRooms(),
    GroupFatigue(), ProfFatigue() {
}

int square(int x) {
    return x * x;
}

fatigue State::partialFatigue(int* schedule) {
    class_time maxTime = 0;
    for (class_time time = ClassesPerDay; time > 0; time--) {
        if (schedule[time] != 0) {
            maxTime = time;
            break;
        }
    }
    if (maxTime == 0) {
        return 0;
    }
    class_time minTime = 0;
    for (class_time time = 1; time <= ClassesPerDay; time++) {
        if (schedule[time] != 0) {
            minTime = time;
            break;
        }
    }
    return square(2 + maxTime - minTime + 1);
}

fatigue State::groupFatigue(group group, day day) {
    return partialFatigue(GroupSchedule[group][day]);
}

fatigue State::profFatigue(prof prof, day day) {
    return partialFatigue(ProfSchedule[prof][day]);
}

struct xor_pair_hash {
    template <class T1, class T2>
    std::size_t operator()(const std::pair<T1,T2> &p) const {
        auto h1 = std::hash<T1>{}(p.first);
        auto h2 = std::hash<T2>{}(p.second);
        return h1 ^ h2;
    }
};

std::unique_ptr<State> SolveNaive(const std::shared_ptr<Problem> problem) {
    auto state = std::make_unique<State>(problem);

    std::unordered_map<std::pair<group, prof>, int, xor_pair_hash> classesToSchedule;
    for (group group = 1; group <= problem->NumGroups; group++) {
        for (prof prof = 1; prof <= problem->NumProfs; prof++) {
            if (problem->NumClasses[group][prof] == 0) {
                continue;
            }
            auto groupAndProf = std::make_pair(group, prof);
            classesToSchedule[groupAndProf] = problem->NumClasses[group][prof];
        }
    }
    for (day day = 1; day <= DaysPerWeek; day++) {
        for (class_time time = 1; time <= ClassesPerDay; time++) {
            state->NumFreeRooms[day][time] = problem->NumRooms;
            std::vector<bool> groupIsBusy(problem->NumGroups + 1);
            std::vector<bool> profIsBusy(problem->NumProfs + 1);
            for (auto it = classesToSchedule.begin(); it != classesToSchedule.end();) {
                auto& groupAndProf = it->first;
                auto& group = groupAndProf.first;
                auto& prof = groupAndProf.second;
                if (groupIsBusy[group] || profIsBusy[prof]) {
                    ++it;
                    continue;
                }
                if (state->NumFreeRooms[day][time] == 0) {
                    break;
                }
                state->GroupSchedule[group][day][time] = prof;
                state->ProfSchedule[prof][day][time] = group;
                groupIsBusy[group] = true;
                profIsBusy[prof] = true;
                state->NumFreeRooms[day][time]--;
                it->second--;
                if (it->second == 0) {
                    it = classesToSchedule.erase(it);
                } else {
                    ++it;
                }
            }
        }
    }

    state->Fatigue = 0;
    for (day day = 1; day <= DaysPerWeek; day++) {
        for (group group = 1; group <= problem->NumGroups; group++) {
            state->GroupFatigue[group][day] = state->groupFatigue(group, day);
            state->Fatigue += state->GroupFatigue[group][day];
        }
        for (prof prof = 1; prof <= problem->NumProfs; prof++) {
            state->ProfFatigue[prof][day] = state->profFatigue(prof, day);
            state->Fatigue += state->ProfFatigue[prof][day];
        }
    }
    return state;
}

int Random(int n) {
    // Note: This generator is not uniform but is probably faster.
    return std::rand() % n;
}

float RandomFloat() {
    return static_cast<float>(std::rand()) / RAND_MAX;
}

bool ShouldAccept(fatigue delta, double progress) {
    if (delta <= 0) {
        return true;
    }
    if (progress < 0.5 || progress > 0.7) {
        return false;
    }
    float temperature = 0.5;
    auto p = exp(static_cast<float>(-delta) / temperature);
    return p >= RandomFloat();
}

std::unique_ptr<Solution> Solve(const std::shared_ptr<Problem> problem) {
    typedef std::chrono::steady_clock clock;
    auto start = clock::now();
    std::ofstream log("./out/log.tsv");
    log << std::fixed << std::setprecision(4);
    log << "Iteration\tFatigue\n";
    auto state = SolveNaive(problem);
    auto bestSolution = std::make_unique<Solution>(*state);
    for (int i = 0; ; i++) {
        auto elapsed = clock::now() - start;
        auto time_left = time_limit - elapsed;
        if (time_left <= clock::duration::zero()) {
            break;
        }
        if (i % 100 == 0) {
            log << i << "\t" << state->Fatigue << "\n";
        }
        for (int t = 0; t < 10; t++) {
            // Generate a swap.
            auto d1 = 1 + Random(DaysPerWeek);
            auto c1 = 1 + Random(ClassesPerDay);
            auto p = 1 + Random(problem->NumProfs);
            auto g = state->ProfSchedule[p][d1][c1];
            if (g == 0) {
                continue;
            }
            auto d2 = 1 + Random(DaysPerWeek);
            auto c2 = 1 + Random(ClassesPerDay);
            if (state->NumFreeRooms[d2][c2] == 0 ||
                    state->ProfSchedule[p][d2][c2] != 0 ||
                    state->GroupSchedule[g][d2][c2] != 0) {
                continue;
            }

            if (1 < c1 && c1 < ClassesPerDay) {
                auto groupWillHaveEmptySlot =
                    state->GroupSchedule[g][d1][c1-1] != 0 &&
                    state->GroupSchedule[g][d1][c1+1] != 0;
                if (groupWillHaveEmptySlot) {
                    continue;
                }
                auto profWillHaveEmptySlot =
                    state->ProfSchedule[p][d1][c1-1] != 0 &&
                    state->ProfSchedule[p][d1][c1+1] != 0;
                if (profWillHaveEmptySlot) {
                    continue;
                }
            }

            auto prevFatigue = state->Fatigue;
            auto prevGroupFatigue1 = state->GroupFatigue[g][d1];
            auto prevGroupFatigue2 = state->GroupFatigue[g][d2];
            auto prevProfFatigue1 = state->ProfFatigue[p][d1];
            auto prevProfFatigue2 = state->ProfFatigue[p][d2];

            // Apply swap.
            state->Fatigue -= state->GroupFatigue[g][d1];
            state->Fatigue -= state->ProfFatigue[p][d1];
            if (d2 != d1) {
                state->Fatigue -= state->GroupFatigue[g][d2];
                state->Fatigue -= state->ProfFatigue[p][d2];
            }
            state->NumFreeRooms[d1][c1]++;
            state->NumFreeRooms[d2][c2]--;
            state->GroupSchedule[g][d1][c1] = 0;
            state->GroupSchedule[g][d2][c2] = p;
            state->ProfSchedule[p][d1][c1] = 0;
            state->ProfSchedule[p][d2][c2] = g;
            state->GroupFatigue[g][d1] = state->groupFatigue(g, d1);
            state->ProfFatigue[p][d1] = state->profFatigue(p, d1);
            state->Fatigue += state->GroupFatigue[g][d1];
            state->Fatigue += state->ProfFatigue[p][d1];
            if (d2 != d1) {
                state->GroupFatigue[g][d2] = state->groupFatigue(g, d2);
                state->ProfFatigue[p][d2] = state->profFatigue(p, d2);
                state->Fatigue += state->GroupFatigue[g][d2];
                state->Fatigue += state->ProfFatigue[p][d2];
            }

            if (state->Fatigue <= prevFatigue) {
                // Accept swap.
                if (state->Fatigue < bestSolution->Fatigue) {
                    bestSolution = std::make_unique<Solution>(*state);
                }
            } else {
                // Revert swap.
                state->NumFreeRooms[d1][c1]--;
                state->NumFreeRooms[d2][c2]++;
                state->GroupSchedule[g][d2][c2] = 0;
                state->GroupSchedule[g][d1][c1] = p;
                state->ProfSchedule[p][d2][c2] = 0;
                state->ProfSchedule[p][d1][c1] = g;
                state->Fatigue = prevFatigue;
                state->GroupFatigue[g][d1] = prevGroupFatigue1;
                state->ProfFatigue[p][d1] = prevProfFatigue1;
                if (d2 != d1) {
                    state->GroupFatigue[g][d2] = prevGroupFatigue2;
                    state->ProfFatigue[p][d2] = prevProfFatigue2;
                }
            }

            break;
        }
    }
    log.close();
    return bestSolution;
}

std::istream& operator>>(std::istream& in, Problem& problem) {
    in >> problem.NumGroups;
    in >> problem.NumProfs;
    in >> problem.NumRooms;
    for (group group = 1; group <= problem.NumGroups; group++) {
        for (prof prof = 1; prof <= problem.NumProfs; prof++) {
            in >> problem.NumClasses[group][prof];
        }
    }
    return in;
}

std::ostream& operator<<(std::ostream& out, Solution& solution) {
    out << solution.Fatigue << "\n";
    auto problem = solution.problem;
    for (group group = 1; group <= problem->NumGroups; group++) {
        out << "\n";
        for (class_time time = 1; time <= ClassesPerDay; time++) {
            for (day day = 1; day <= DaysPerWeek; day++) {
                if (day != 1) {
                    out << " ";
                }
                out << solution.GroupSchedule[group][day][time];
            }
            out << "\n";
        }
    }
    return out;
}

}  // namespace

int main() {
    std::ios_base::sync_with_stdio(false);
    auto problem = std::make_shared<Problem>();
    std::cin >> *problem;
    auto solution = Solve(problem);
    std::cout << *solution;
    return 0;
}
