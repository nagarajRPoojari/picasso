load("@rules_cc//cc:defs.bzl", "cc_library")

cc_library(
    name = "libffi",
    hdrs = glob(["include/**/*.h"]),
    srcs = glob(["lib/aarch64-linux-gnu/libffi.so*"]),
    includes = ["include"],
    visibility = ["//visibility:public"],
)