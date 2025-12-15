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
            "bin/irgen gen %s .niyama/out.ll", file);
        if (system(cmd) != 0) die("irgen");

        if (system("llvm-as-16 .niyama/out.ll -o .niyama/out.bc") != 0)
            die("llvm-as");

        if (system("llc-16 -filetype=obj .niyama/out.bc -o .niyama/out.o") != 0)
            die("llc");

        if (system("cc .niyama/out.o ./bin/libruntime.a -o .niyama/a.out -luring -lpthread -lm") != 0)
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