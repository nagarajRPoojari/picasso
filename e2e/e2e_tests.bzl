load("@rules_shell//shell:sh_test.bzl", "sh_test")

def e2e_tests(name, runner, picasso):
    test_dirs = native.glob(["**/TEST_DIR"])
    tests = []

    for marker in test_dirs:
        test_dir = marker[:-len("/TEST_DIR")]
        test_name = test_dir.replace("/", ".")
        test_path = test_dir.split("/")[-1]
        filegroup_name = test_name + "_files"

        native.filegroup(
            name = filegroup_name,
            srcs = native.glob([test_dir + "/**"]),
            visibility = ["//visibility:private"],
        )

        target_name = test_name + ".test"

        sh_test(
            name = target_name,
            srcs = ["run_single_e2e.sh"],
            args = [
                "$(rootpath %s)" % runner,
                "$(rootpath %s)" % picasso,
                test_path,
                "--deps",
                "$(locations //irgen)",
                "$(locations //:runtime_lib)",
                "--files",
                "$(locations :%s)" % filegroup_name,
            ],
            data = [
                runner,
                picasso,
                "//irgen",
                "//:runtime_lib",
                ":%s" % filegroup_name,
                "//:runtime_headers",
                "//libs:ffi_libs", 
            ],
            timeout = "long",
            visibility = ["//visibility:public"],
        )
        tests.append(":" + target_name)

    native.test_suite(
        name = name,
        tests = tests,
        visibility = ["//visibility:public"],
    )