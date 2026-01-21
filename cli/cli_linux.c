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

static void get_tool_root(char *out, size_t sz) {
    const char *runfiles = getenv("RUNFILES_DIR");
    if (runfiles) {
        if (snprintf(out, sz, "%s/_main", runfiles) >= (int)sz)
            die("toolRoot path too long");
    } else {
        if (getcwd(out, sz) == NULL) {
            perror("getcwd");
            exit(1);
        }
    }
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

static void generate_ffi_irs_from_root( const char *ffiRoot, const char *tmpDir, const char *ffiObjDir, const char *runtimeIncDir, int gen_objects) {
    DIR *d = opendir(ffiRoot);
    if (!d && !runtimeIncDir) return;
    if (!d) die("opendir ffiRoot");

    struct dirent *ent;
    while ((ent = readdir(d)) != NULL) {
        if (ent->d_name[0] == '.') continue;

        char tuDir[PATH_MAX];
        if (snprintf(tuDir, sizeof(tuDir), "%s/%s", ffiRoot, ent->d_name) >= (int)sizeof(tuDir))
            die("tuDir path too long");

        struct stat st;
        if (stat(tuDir, &st) != 0 || !S_ISDIR(st.st_mode)) continue;

        char stub[PATH_MAX];
        if (snprintf(stub, sizeof(stub), "%s/ffi_stub.c", tuDir) >= (int)sizeof(stub))
            die("stub path too long");
            
        if (access(stub, R_OK) != 0) {
            if (snprintf(stub, sizeof(stub), "%s/ffi_stub_linux.c", tuDir) >= (int)sizeof(stub))
                die("stub path too long");

            if(access(stub, R_OK) != 0)  {
                fprintf(stderr, "warning: %s missing ffi_stub.c, skipping\n", tuDir);
                continue;
            }
        }

        char outLL[PATH_MAX];
        if (snprintf(outLL, sizeof(outLL), "%s/%s.ll", tmpDir, ent->d_name) >= (int)sizeof(outLL))
            die("outLL path too long");

        /* clang -S -emit-llvm */
        char *clang_ll[20];
        int i = 0;

        clang_ll[i++] = "clang";
        clang_ll[i++] = "-S";
        clang_ll[i++] = "-emit-llvm";
        clang_ll[i++] = "-g0";

        /* runtime headers first */
        if (runtimeIncDir) {
            clang_ll[i++] = "-I";
            clang_ll[i++] = (char *)runtimeIncDir;
        }

        /* TU-local headers */
        clang_ll[i++] = "-I";
        clang_ll[i++] = (char *)tuDir;

        clang_ll[i++] = (char *)stub;
        clang_ll[i++] = "-o";
        clang_ll[i++] = outLL;
        clang_ll[i++] = NULL;

        run_cmd(clang_ll);

        if (!gen_objects) continue;

        /* compile other .c → native .o */
        DIR *cd = opendir(tuDir);
        if (!cd) die("opendir tuDir");

        struct dirent *cent;
        while ((cent = readdir(cd)) != NULL) {
            if (!strstr(cent->d_name, ".c")) continue;
            if (!strcmp(cent->d_name, "ffi_stub.c")) continue;
            if (!strcmp(cent->d_name, "ffi_stub_linux.c")) continue;

            char cfile[PATH_MAX];
            if (snprintf(cfile, sizeof(cfile), "%s/%s", tuDir, cent->d_name) >= (int)sizeof(cfile))
                die("cfile path too long");

            char obj[PATH_MAX];
            if (snprintf(obj, sizeof(obj), "%s/%s_%s.o", ffiObjDir, ent->d_name, cent->d_name) >= (int)sizeof(obj))
                die("obj path too long");

            char *cc_obj[16];
            i = 0;
            cc_obj[i++] = "cc";
            cc_obj[i++] = "-c";
            cc_obj[i++] = "-fPIC";
            cc_obj[i++] = "-I";
            cc_obj[i++] = (char *)tuDir;
            if (runtimeIncDir) {
                cc_obj[i++] = "-I";
                cc_obj[i++] = (char *)runtimeIncDir;
            }
            cc_obj[i++] = cfile;
            cc_obj[i++] = "-o";
            cc_obj[i++] = obj;
            cc_obj[i++] = NULL;

            run_cmd(cc_obj);
        }
        closedir(cd);
    }
    closedir(d);
}

static void generate_ffi_irs(const char *dir, const char *buildDir) {
    char ffiRoot[PATH_MAX];
    snprintf(ffiRoot, sizeof(ffiRoot), "%s/c/ffi", dir);

    char tmpDir[PATH_MAX];
    snprintf(tmpDir, sizeof(tmpDir), "%s/tmp", buildDir);

    char ffiObjDir[PATH_MAX];
    snprintf(ffiObjDir, sizeof(ffiObjDir), "%s/tmp/ffi-obj", buildDir);

    run_cmd((char *[]){"mkdir", "-p", tmpDir, NULL});
    run_cmd((char *[]){"mkdir", "-p", ffiObjDir, NULL});

    generate_ffi_irs_from_root(
        ffiRoot,
        tmpDir,
        ffiObjDir,
        NULL,
        1
    );
}

/* main */
int main(int argc, char **argv) {
    if (argc < 2) {
        fprintf(stderr, "usage: picasso <build|exec> <project-dir>\n");
        return 1;
    }

    if (!strcmp(argv[1], "build")) {
        if (argc != 3) {
            fprintf(stderr, "picasso build <project root dir>\n");
            return 1;
        }

        const char *dir = argv[2];

        char buildDir[PATH_MAX];
        snprintf(buildDir, sizeof(buildDir), "%s/build", dir);

        run_cmd((char *[]){"mkdir", "-p", buildDir, NULL});

        /* project FFI IR + objects */
        generate_ffi_irs(dir, buildDir);

        /* stdlib FFI IRs */
        char toolRoot[PATH_MAX];
        get_tool_root(toolRoot, sizeof(toolRoot));

        char libsDir[PATH_MAX];
        snprintf(libsDir, sizeof(libsDir), "%s/libs", toolRoot);

        char runtimeHdrs[PATH_MAX];
        snprintf(runtimeHdrs, sizeof(runtimeHdrs), "%s/runtime/headers", toolRoot);

        char tmpDir[PATH_MAX];
        snprintf(tmpDir, sizeof(tmpDir), "%s/build/tmp", dir);

        generate_ffi_irs_from_root(
            libsDir,
            tmpDir,
            NULL,
            runtimeHdrs,
            0
        );

        /* language IR generation */
        char *irgen[] = { IRGEN_BIN, "gen", (char *)dir, NULL };
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
            "set -e; "
            "OBJS=\"$1\"/*.o; "
            "FFI_OBJS=\"\"; "
            "if [ -d \"$1/tmp/ffi-obj\" ] && ls \"$1/tmp/ffi-obj\"/*.o >/dev/null 2>&1; then "
            "  FFI_OBJS=\"$1/tmp/ffi-obj\"/*.o; "
            "fi; "
            "cc $OBJS $FFI_OBJS "
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
            fprintf(stderr, "picasso exec <project root dir>\n");
            return 1;
        }

        char exe[PATH_MAX];
        snprintf(exe, sizeof(exe), "%s/build/a.out", argv[2]);
        execl(exe, exe, (char *)NULL);
        die("exec");
    }

    fprintf(stderr, "unknown command\n");
    return 1;
}
