#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/wait.h>
#include <sys/stat.h>
#include <string.h>

static void run_cmd(const char *bin, const char *cmd, const char *dir) {
    pid_t pid = fork();
    if (pid == 0) {
        if (chdir(dir) != 0) { perror("chdir"); _exit(1); }
        execl(bin, bin, cmd, ".", NULL);
        perror("exec failed");
        _exit(1);
    }

    int status;
    waitpid(pid, &status, 0);
    if (!WIFEXITED(status) || WEXITSTATUS(status) != 0) {
        fprintf(stderr, "ERROR: Command failed\n");
        exit(1);
    }
}

int main(int argc, char **argv) {
    if (argc < 3) return 2;
    run_cmd(argv[1], "build", argv[2]);
    run_cmd(argv[1], "exec", argv[2]);
    printf("\n✓ All steps passed!\n");
    return 0;
}