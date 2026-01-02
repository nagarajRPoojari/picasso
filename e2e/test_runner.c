#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/wait.h>
#include <sys/stat.h>
#include <string.h>
#include <string.h>
#include <ctype.h>

#define OUTPUT_BUFFER_SIZE 4096

typedef struct {
    char section[64];
    char key[64];
    char value[256];
} ini_entry_t;

static char *trim(char *s) {
    while (isspace((unsigned char)*s)) s++;
    if (*s == 0) return s;

    char *end = s + strlen(s) - 1;
    while (end > s && isspace((unsigned char)*end)) end--;
    end[1] = '\0';
    return s;
}

int parse_ini(const char *filename, ini_entry_t *entries, int max_entries) {
    FILE *fp = fopen(filename, "r");
    if (!fp) return -1;

    char line[512];
    char current_section[64] = "";
    int count = 0;

    while (fgets(line, sizeof(line), fp)) {
        char *s = trim(line);

        // Skip empty lines & comments
        if (*s == '\0' || *s == '#' || *s == ';')
            continue;

        // Section
        if (*s == '[') {
            char *end = strchr(s, ']');
            if (!end) continue;

            *end = '\0';
            strncpy(current_section, s + 1, sizeof(current_section));
            continue;
        }

        // Key = Value
        char *eq = strchr(s, '=');
        if (!eq) continue;

        *eq = '\0';
        char *key = trim(s);
        char *value = trim(eq + 1);

        if (count < max_entries) {
            strncpy(entries[count].section, current_section, 64);
            strncpy(entries[count].key, key, 64);
            strncpy(entries[count].value, value, 256);
            count++;
        }
    }

    fclose(fp);
    return count;
}

static const char *ini_get(ini_entry_t *entries, int n, const char *section, const char *key) {
    for (int i = 0; i < n; i++) {
        if (strcmp(entries[i].section, section) == 0 &&
            strcmp(entries[i].key, key) == 0) {
            return entries[i].value;
        }
    }
    return NULL;
}


static int run_cmd(const char *bin, const char *cmd, const char *dir) {
    pid_t pid = fork();
    if (pid == 0) {
        if (chdir(dir) != 0) {
            perror("chdir");
            _exit(127);
        }
        execl(bin, bin, cmd, ".", NULL);
        perror("exec failed");
        _exit(127);
    }

    int status;
    waitpid(pid, &status, 0);

    if (WIFEXITED(status))
        return WEXITSTATUS(status) == 0;

    return 0;
}


static int expect_pass(const char *v) {
    return v && strcmp(v, "pass") == 0;
}


static int read_file(const char *path, char *buf, size_t max) {
    FILE *f = fopen(path, "r");
    if (!f) return -1;

    size_t n = fread(buf, 1, max - 1, f);
    buf[n] = '\0';
    fclose(f);
    return 0;
}

static int run_cmd_capture( const char *bin, const char *cmd, const char *dir, char *out, size_t out_sz ) {
    int pipefd[2];
    if (pipe(pipefd) != 0)
        return 0;

    pid_t pid = fork();
    if (pid == 0) {
        close(pipefd[0]);
        dup2(pipefd[1], STDOUT_FILENO);
        close(pipefd[1]);

        if (chdir(dir) != 0)
            _exit(127);

        execl(bin, bin, cmd, ".", NULL);
        _exit(127);
    }

    close(pipefd[1]);

    size_t total = 0;
    ssize_t n;

    while ((n = read(pipefd[0], out + total, out_sz - 1 - total)) > 0) {
        total += n;
        if (total >= out_sz - 1)
            break;
    }

    out[total] = '\0';
    close(pipefd[0]);

    int status;
    waitpid(pid, &status, 0);

    if (WIFEXITED(status))
        return WEXITSTATUS(status) == 0;

    return 0;
}


static int is_yes(const char *v) {
    return v && strcmp(v, "yes") == 0;
}


int main(int argc, char **argv) {
    if (argc < 3) {
        fprintf(stderr, "usage: %s <runner> <dir>\n", argv[0]);
        return 2;
    }

    char config_path[512];
    snprintf(config_path, sizeof(config_path),
             "%s/config.ini", argv[2]);

    ini_entry_t entries[128];
    int n = parse_ini(config_path, entries, 128);
    if (n < 0) {
        perror("parse_ini");
        return 1;
    }

    const char *build_exp = ini_get(entries, n, "build", "status");
    const char *exec_exp  = ini_get(entries, n, "exec", "status");
    const char *verify_out = ini_get(entries, n, "output", "verify");


    printf("→ running build\n");
    int build_ok = run_cmd(argv[1], "build", argv[2]);

    if (build_ok != expect_pass(build_exp)) {
        fprintf(stderr, "ERROR: build expected %s but %s\n", build_exp, build_ok ? "passed" : "failed");
        return 1;
    }

    printf("→ running exec\n");

    char exec_out[OUTPUT_BUFFER_SIZE];
    int exec_ok = run_cmd_capture(argv[1], "exec", argv[2], exec_out, sizeof(exec_out));

    if (exec_ok != expect_pass(exec_exp)) {
        fprintf(stderr, "ERROR: exec expected %s but %s\n", exec_exp, exec_ok ? "passed" : "failed");
        return 1;
    }

    if (is_yes(verify_out)) {
        char expected[OUTPUT_BUFFER_SIZE];
        char out_path[512];

        snprintf(out_path, sizeof(out_path), "%s/output.txt", argv[2]);

        if (read_file(out_path, expected, sizeof(expected)) != 0) {
            perror("read output.txt");
            return 1;
        }

        if (strcmp(exec_out, expected) != 0) {
            fprintf(stderr, "ERROR: stdout mismatch\n");
            fprintf(stderr, "---- expected ----\n%s\n", expected);
            fprintf(stderr, "---- got ----\n%s\n", exec_out);
            return 1;
        }
    }
    return 0;
}