%DirectoryReader = type { i64 }
%Math = type { i64 }

@.str1 = private global [22 x i8] c"assigining - %f to  \0A\00"
@.str2 = private global [9 x i8] c" - %s  \0A\00"
@.str3 = private global [2 x i8] c"n\00"
@.str4 = private global [15 x i8] c"instance - %s\0A\00"
@.str5 = private global [16 x i8] c"DirectoryReader\00"
@.str6 = private global [22 x i8] c"assigining - %f to  \0A\00"
@.str7 = private global [9 x i8] c" - %s  \0A\00"
@.str8 = private global [2 x i8] c"a\00"
@.str9 = private global [22 x i8] c"assigining - %f to  \0A\00"
@.str10 = private global [9 x i8] c" - %s  \0A\00"
@.str11 = private global [2 x i8] c"z\00"
@.str12 = private global [18 x i8] c"return z = %f %d\0A\00"

define i64 @main.DirectoryReader.sum(i64 %x, %DirectoryReader %this) {
entry:
	%0 = alloca i64
	store i64 %x, i64* %0
	%1 = alloca %DirectoryReader
	%2 = getelementptr %DirectoryReader, %DirectoryReader* %1, i32 0, i32 0
	%3 = load i64, i64* %2
	%4 = getelementptr [22 x i8], [22 x i8]* @.str1, i64 0, i64 0
	%5 = call i32 (i8*, ...) @printf(i8* %4, i64 %3)
	%6 = getelementptr [9 x i8], [9 x i8]* @.str2, i64 0, i64 0
	%7 = getelementptr [2 x i8], [2 x i8]* @.str3, i64 0, i64 0
	%8 = call i32 (i8*, ...) @printf(i8* %6, i8* %7)
	%9 = getelementptr %DirectoryReader, %DirectoryReader* %1, i32 0, i32 0
	%10 = load i64, i64* %9
	%11 = load i64, i64* %0
	%12 = sitofp i64 %10 to double
	%13 = sitofp i64 %11 to double
	%14 = fadd double %12, %13
	%15 = alloca double
	store double %14, double* %15
	%16 = load double, double* %15
	%17 = fptosi double %16 to i64
	ret i64 %17
}

declare i32 @printf(i8* %0, ...)

define i32 @main() {
entry:
	%0 = alloca %DirectoryReader
	%1 = alloca double
	store double 112.0, double* %1
	%2 = load double, double* %1
	%3 = bitcast double %2 to i64
	%4 = getelementptr %DirectoryReader, %DirectoryReader* %0, i32 0, i32 0
	store i64 %3, i64* %4
	%5 = getelementptr [15 x i8], [15 x i8]* @.str4, i64 0, i64 0
	%6 = getelementptr [16 x i8], [16 x i8]* @.str5, i64 0, i64 0
	%7 = call i32 (i8*, ...) @printf(i8* %5, i8* %6)
	%8 = load %DirectoryReader, %DirectoryReader* %0
	%9 = getelementptr [22 x i8], [22 x i8]* @.str6, i64 0, i64 0
	%10 = call i32 (i8*, ...) @printf(i8* %9, %DirectoryReader %8)
	%11 = getelementptr [9 x i8], [9 x i8]* @.str7, i64 0, i64 0
	%12 = getelementptr [2 x i8], [2 x i8]* @.str8, i64 0, i64 0
	%13 = call i32 (i8*, ...) @printf(i8* %11, i8* %12)
	%14 = alloca double
	store double 10.0, double* %14
	%15 = load double, double* %14
	%16 = fptosi double %15 to i64
	%17 = load %DirectoryReader, %DirectoryReader* %0
	%18 = call i64 @main.DirectoryReader.sum(i64 %16, %DirectoryReader %17)
	%19 = alloca i64
	store i64 %18, i64* %19
	%20 = load i64, i64* %19
	%21 = getelementptr [22 x i8], [22 x i8]* @.str9, i64 0, i64 0
	%22 = call i32 (i8*, ...) @printf(i8* %21, i64 %20)
	%23 = getelementptr [9 x i8], [9 x i8]* @.str10, i64 0, i64 0
	%24 = getelementptr [2 x i8], [2 x i8]* @.str11, i64 0, i64 0
	%25 = call i32 (i8*, ...) @printf(i8* %23, i8* %24)
	%26 = alloca double
	store double 0.0, double* %26
	%27 = load double, double* %26
	%28 = fptosi double %27 to i32
	%29 = load i64, i64* %19
	%30 = getelementptr [18 x i8], [18 x i8]* @.str12, i64 0, i64 0
	%31 = call i32 (i8*, ...) @printf(i8* %30, i64 %29, i64 %29)
	ret i32 %28
}
