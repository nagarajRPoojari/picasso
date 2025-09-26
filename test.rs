			import io;
			import types;
			fn main(): int32 {
				say a: float16 = 65504;
				say b: int16 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}