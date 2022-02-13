// Code generated by gioui.org/cpu/cmd/compile DO NOT EDIT.

struct backdrop_descriptor_set_layout {
	struct buffer_descriptor binding0;
	struct buffer_descriptor binding1;
};

extern coroutine backdrop_coroutine_begin(struct program_data *data,
	int32_t workgroupX, int32_t workgroupY, int32_t workgroupZ,
	void *workgroupMemory,
	int32_t firstSubgroup,
	int32_t subgroupCount) ATTR_HIDDEN;

extern bool backdrop_coroutine_await(coroutine r, yield_result *res) ATTR_HIDDEN;
extern void backdrop_coroutine_destroy(coroutine r) ATTR_HIDDEN;

extern const struct program_info backdrop_program_info ATTR_HIDDEN;
