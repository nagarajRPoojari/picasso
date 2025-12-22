#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

static void die(const char *msg) {
    perror(msg);
    exit(1);
}

int main(int argc, char **argv) {
    if (argc < 2) {
        fprintf(stderr, "usage: niyama <build|exec> [file]\n");
        return 1;
    }

    if (!strcmp(argv[1], "build")) {
        if (argc != 3) {
            fprintf(stderr, "niyama build <project root dir>\n");
            return 1;
        }

        const char *dir = argv[2];

        char buildDir[1024];
        snprintf(buildDir, sizeof(buildDir), "%s/build", dir);

        char cmd[4096];

        /* create build dir */
        snprintf(cmd, sizeof(cmd), "mkdir -p \"%s\"", buildDir);
        if (system(cmd) != 0) die("mkdir build");

        /* generate LLVM IR (.ll) */
        snprintf(cmd, sizeof(cmd),
                "bin/irgen gen \"%s\"",
                dir);
        if (system(cmd) != 0) die("irgen");

        /* .ll -> .bc */
        snprintf(cmd, sizeof(cmd),
            "set -e; "
            "for f in \"%s\"/*.ll; do "
            "  base=$(basename \"$f\" .ll); "
            "  llvm-as-16 \"$f\" -o \"%s/$base.bc\"; "
            "done",
            buildDir, buildDir);
        if (system(cmd) != 0) die("llvm-as");

        /* .bc -> .o */
        snprintf(cmd, sizeof(cmd),
            "set -e; "
            "for f in \"%s\"/*.bc; do "
            "  base=$(basename \"$f\" .bc); "
            "  llc-16 -filetype=obj \"$f\" -o \"%s/$base.o\"; "
            "done",
            buildDir, buildDir);
        if (system(cmd) != 0) die("llc");

        /* link */
        snprintf(cmd, sizeof(cmd),
            "cc \"%s\"/*.o bin/libruntime.a "
            "-o \"%s\"/a.out "
            "-luring -lunwind-aarch64 -lpthread -lm",
            buildDir, buildDir);
        if (system(cmd) != 0) die("link");

        return 0;
    }



    if (!strcmp(argv[1], "exec")) {
        if (argc != 3) {
            fprintf(stderr, "niyama exec <project root dir>\n");
            return 1;
        }

        const char *dir = argv[2];

        char exe[1024];
        snprintf(exe, sizeof(exe), "%s/build/a.out", dir);

        execl(exe, exe, (char *)NULL);
        die("exec");
    }


    fprintf(stderr, "unknown command\n");
    return 1;
}