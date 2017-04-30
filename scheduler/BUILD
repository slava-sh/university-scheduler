cc_library(
    name = "scheduler",
    srcs = glob(["*.cpp"], exclude=["main.cpp"]),
    hdrs = glob(["*.h"]),
)

cc_binary(
    name = "main",
    srcs = ["main.cpp"],
    deps = [":scheduler"],
)
