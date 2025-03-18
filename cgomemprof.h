#include <stddef.h>
#include "backtrace.h"

#ifndef CGOMEMPROF_H
#define CGOMEMPROF_H

// EnableMemoryProfiling enables memory profiling.
int EnableMemoryProfiling();

// DisableMemoryProfiling disables memory profiling.
int DisableMemoryProfiling();

// DumpMemoryProfileIntoFile dumps the memory profile into the specified
// filename.
int DumpMemoryProfileIntoFile(const char *filename);

char*  GetSymbol(uintptr_t addr);

// implemet by jemalloc
extern int mallctl(const char *name, void *oldp, size_t *oldlenp, void *newp,
                   size_t newlen);

int IsJemallocProfEnabled();

#endif // CGOMEMPROF_H