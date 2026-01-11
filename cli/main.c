#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <dirent.h>
#include <sys/stat.h>
#include <sys/wait.h>
#include <limits.h>
#include <errno.h>

static void die(const char *msg) {
    perror(msg);
    exit(1);
}

static void run_cmd(char *const argv[]) {
    pid_t pid = fork();
    if (pid < 0)
        die("fork");

    if (pid == 0) {
        execvp(argv[0], argv);
        perror(argv[0]);
        _exit(127);
    }

    int status;
    if (waitpid(pid, &status, 0) < 0)
        die("waitpid");

    if (!WIFEXITED(status) || WEXITSTATUS(status) != 0) {
        fprintf(stderr, "command failed: %s\n", argv[0]);
        exit(1);
    }
}

/* FFI IR + native object generation */
static void generate_ffi_irs(const char *dir, const char *buildDir) {
    char ffiRoot[PATH_MAX];
    if (snprintf(ffiRoot, sizeof(ffiRoot), "%s/c/ffi", dir) >= (int)sizeof(ffiRoot))
        die("ffiRoot path too long");

    char tmpDir[PATH_MAX];
    if (snprintf(tmpDir, sizeof(tmpDir), "%s/tmp", buildDir) >= (int)sizeof(tmpDir))
        die("tmpDir path too long");

    char ffiObjDir[PATH_MAX];
    if (snprintf(ffiObjDir, sizeof(ffiObjDir),
                 "%s/tmp/ffi-obj", buildDir) >= (int)sizeof(ffiObjDir))
        die("ffiObjDir path too long");

    /* mkdir -p build/tmp build/tmp/ffi-obj */
    char *mkdir_tmp[] = { "mkdir", "-p", tmpDir, NULL };
    run_cmd(mkdir_tmp);

    char *mkdir_obj[] = { "mkdir", "-p", ffiObjDir, NULL };
    run_cmd(mkdir_obj);

    DIR *d = opendir(ffiRoot);
    if (!d) die("opendir c/ffi");

    struct dirent *ent;
    while ((ent = readdir(d)) != NULL) {
        if (ent->d_name[0] == '.')
            continue;

        char tuDir[PATH_MAX];
        if (snprintf(tuDir, sizeof(tuDir),
                     "%s/%s", ffiRoot, ent->d_name) >= (int)sizeof(tuDir))
            die("tuDir path too long");

        struct stat st;
        if (stat(tuDir, &st) != 0 || !S_ISDIR(st.st_mode))
            continue;

        /* ffi_stub.c → LLVM IR */

        char stub[PATH_MAX];
        if (snprintf(stub, sizeof(stub),
                     "%s/ffi_stub.c", tuDir) >= (int)sizeof(stub))
            die("stub path too long");

        if (access(stub, R_OK) != 0) {
            fprintf(stderr, "warning: %s missing ffi_stub.c, skipping TU\n", tuDir);
            continue;
        }

        char outLL[PATH_MAX];
        if (snprintf(outLL, sizeof(outLL),
                     "%s/%s.ll", tmpDir, ent->d_name) >= (int)sizeof(outLL))
            die("outLL path too long");

        printf("FFI IRGEN: %s\n", ent->d_name);

        char *clang_ll[] = {
            "clang",
            "-S",
            "-emit-llvm",
            "-g0",
            "-I", tuDir,
            stub,
            "-o", outLL,
            NULL
        };
        run_cmd(clang_ll);

        /* compile all other .c → native .o */
        DIR *cd = opendir(tuDir);
        if (!cd)
            die("opendir tuDir");

        struct dirent *cent;
        while ((cent = readdir(cd)) != NULL) {
            if (!strstr(cent->d_name, ".c"))
                continue;

            if (!strcmp(cent->d_name, "ffi_stub.c"))
                continue;

            char cfile[PATH_MAX];
            if (snprintf(cfile, sizeof(cfile),
                         "%s/%s", tuDir, cent->d_name) >= (int)sizeof(cfile))
                die("cfile path too long");

            char obj[PATH_MAX];
            if (snprintf(obj, sizeof(obj),
                         "%s/%s_%s.o",
                         ffiObjDir,
                         ent->d_name,
                         cent->d_name) >= (int)sizeof(obj))
                die("obj path too long");

            /* strip ".c" */
            // obj[strlen(obj) - 2] = '\0';

            char *cc_obj[] = {
                "cc",
                "-c",
                "-fPIC",
                "-I", tuDir,
                cfile,
                "-o", obj,
                NULL
            };

            run_cmd(cc_obj);
        }
        closedir(cd);
    }

    closedir(d);
}

/* main */
int main(int argc, char **argv) {
    if (argc < 2) {
        fprintf(stderr, "usage: niyama <build|exec> <project-dir>\n");
        return 1;
    }

    if (!strcmp(argv[1], "build")) {
        if (argc != 3) {
            fprintf(stderr, "niyama build <project root dir>\n");
            return 1;
        }

        const char *dir = argv[2];

        char buildDir[PATH_MAX];
        if (snprintf(buildDir, sizeof(buildDir),
                     "%s/build", dir) >= (int)sizeof(buildDir))
            die("buildDir path too long");

        char *mkdir_build[] = { "mkdir", "-p", buildDir, NULL };
        run_cmd(mkdir_build);

        /* FFI IR + native objects */
        generate_ffi_irs(dir, buildDir);

        /* language IR generation */
        char *irgen[] = {
            IRGEN_BIN,
            "gen",
            (char *)dir,
            NULL
        };
        run_cmd(irgen);

        /* .ll → .bc */
        char *ll_to_bc[] = {
            "sh", "-c",
            "set -e; for f in \"$1\"/*.ll; do "
            "b=$(basename \"$f\" .ll); "
            "llvm-as-16 \"$f\" -o \"$1/$b.bc\"; "
            "done",
            "sh",
            buildDir,
            NULL
        };
        run_cmd(ll_to_bc);

        /* .bc → .o */
        char *bc_to_o[] = {
            "sh", "-c",
            "set -e; for f in \"$1\"/*.bc; do "
            "b=$(basename \"$f\" .bc); "
            "llc-16 -filetype=obj \"$f\" -o \"$1/$b.o\"; "
            "done",
            "sh",
            buildDir,
            NULL
        };
        run_cmd(bc_to_o);

        /* final link */
        char *link[] = {
            "sh", "-c",
            "cc \"$1\"/*.o \"$1\"/tmp/ffi-obj/*.o "
            RUNTIME_LIB_PATH
            " -o \"$1/a.out\" "
            "-rdynamic "
            "-L/usr/lib/aarch64-linux-gnu -L/usr/lib "
            "-lffi -luring -lunwind -lunwind-aarch64 "
            "-lpthread -lm",
            "sh",
            buildDir,
            NULL
        };
        run_cmd(link);

        return 0;
    }

    if (!strcmp(argv[1], "exec")) {
        if (argc != 3) {
            fprintf(stderr, "niyama exec <project root dir>\n");
            return 1;
        }

        char exe[PATH_MAX];
        if (snprintf(exe, sizeof(exe),
                     "%s/build/a.out", argv[2]) >= (int)sizeof(exe))
            die("exe path too long");

        execl(exe, exe, (char *)NULL);
        die("exec");
    }

    fprintf(stderr, "unknown command\n");
    return 1;
}
