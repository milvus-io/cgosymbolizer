//go:build with_jemalloc && linux
// +build with_jemalloc,linux

#include "cgomemprof.h"
#include <stdbool.h>
#include <stdlib.h>
#include <dlfcn.h>
#include "backtrace.h"

extern struct backtrace_state* cgoBacktraceState;

#ifndef RTLD_DEFAULT
#define RTLD_DEFAULT ((void*)0)
#endif

int EnableMemoryProfiling() {
  bool enable = true;
  return mallctl("prof.active", NULL, NULL, &enable, sizeof(enable));
}

int DisableMemoryProfiling() {
  bool enable = false;
  return mallctl("prof.active", NULL, NULL, &enable, sizeof(enable));
}

int DumpMemoryProfileIntoFile(const char *filename) {
  return mallctl("prof.dump", NULL, NULL, &filename, sizeof(const char *));
}

void syminfoCallback(void *data, uintptr_t pc, const char *symname, uintptr_t symval, uintptr_t symsize) {
    char* buf = (char*)data;
    if (symname != NULL) {
        snprintf(buf, 4096, "0x%lx\t%s",pc, symname);
    } else {
        snprintf(buf, 4096, "0x%lx\t??", pc);
    }
}

void errorCallback(void *data, const char *msg, int errnum) {
    fprintf(stderr, "Error: %s (code: %d)\n", msg, errnum);
}

char* GetSymbol(uintptr_t addr) {
    char* buf = malloc(4096);
    backtrace_syminfo(cgoBacktraceState, addr, syminfoCallback, errorCallback, buf);
    return buf;
}

int IsJemallocProfEnabled() {
    bool enabled = false;
    size_t size = sizeof(enabled);
    int code = mallctl("opt.prof", &enabled, &size, NULL, 0);
    if (code != 0) {
        fprintf(stderr, "mallctl failed with code: %d\n", code);
        return -1;
    }
    return enabled ? 1 : 0;
}
