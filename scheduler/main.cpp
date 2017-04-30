#include <fstream>
#include <iostream>
#include <memory>
#include "scheduler/time_limited_solver.h"

int main(int argc, char** argv) {
  std::ios_base::sync_with_stdio(false);

  auto time_limit = std::chrono::seconds(10) - std::chrono::milliseconds(50);
  if (argc > 1 && std::string(argv[1]) == "--fast") {
      time_limit = std::chrono::milliseconds(300);
  }

  auto problem = std::make_shared<scheduler::Problem>();
  std::cin >> *problem;
  auto solution = scheduler::TimeLimitedSolver(time_limit).Solve(problem);
  std::cout << solution;
  return 0;
}
