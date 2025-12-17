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
            fprintf(stderr, "niyama build <file.ini>\n");
            return 1;
        }

        const char *file = argv[2];

        if (system("mkdir -p .niyama") != 0) die("mkdir");

        char cmd[1024];
        snprintf(cmd, sizeof(cmd),
                "bin/irgen gen %s ./bin", file);
        if (system(cmd) != 0) die("irgen");

        if (system("for f in ./bin/*.ll; do "
                "llvm-as-16 \"$f\" -o \"${f%.ll}.bc\" || exit 1; "
                "done") != 0)
            die("llvm-as");

        if (system("for f in ./bin/*.bc; do "
                "llc-16 -filetype=obj \"$f\" -o \"${f%.bc}.o\" || exit 1; "
                "done") != 0)
            die("llc");

        if (system("cc ./bin/*.o ./bin/libruntime.a "
                "-o .niyama/a.out "
                "-luring -lunwind-aarch64 -lpthread -lm") != 0)
            die("link");

        return 0;
    }


    if (!strcmp(argv[1], "exec")) {
        // Run the generated binary
        execl("./.niyama/a.out", "./.niyama/a.out", NULL);
        die("exec");
    }

    fprintf(stderr, "unknown command\n");
    return 1;
}