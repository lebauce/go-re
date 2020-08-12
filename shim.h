#ifndef __BPF_COMPILER_defined
#define __BPF_COMPILER_defined 1

typedef struct bpf_compiler bpf_compiler;

bpf_compiler *new_bpf_compiler(void);
int bpf_compile_to_object_file(bpf_compiler *, const char *, const char *, const char **, char verbose);
const char * bpf_compiler_get_errors(bpf_compiler *compiler);
void delete_bpf_compiler(bpf_compiler *);

#endif