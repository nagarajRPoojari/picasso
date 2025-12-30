def e2e_tests(name, runner, niyama):
    test_dirs = native.glob(["*/TEST_DIR"])
    tests = []

    for marker in test_dirs:
        test_dir = marker[:-len("/TEST_DIR")]
        test_name = test_dir.split("/")[-1]
        filegroup_name = test_name + "_files"

        native.filegroup(
            name = filegroup_name,
            srcs = native.glob([test_dir + "/**"]),
            visibility = ["//visibility:private"],
        )

        target_name = test_name + ".test"

        native.sh_test(
            name = target_name,
            srcs = ["run_single_e2e.sh"],
            args = [
                "$(rootpath %s)" % runner,
                "$(rootpath %s)" % niyama,
                test_name,
                "--deps",
                "$(locations //irgen)",
                "$(locations //:runtime_lib)",
                "--files",
                "$(locations :%s)" % filegroup_name,
            ],
            data = [
                runner,
                niyama,
                "//irgen",
                "//:runtime_lib",
                ":%s" % filegroup_name,
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