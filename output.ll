%DirectoryReader = type { double, %Math }
%Math = type { double }

@.str1 = private global [15 x i8] c"instance - %s\0A\00"
@.str2 = private global [5 x i8] c"Math\00"
@.str3 = private global [15 x i8] c"instance - %s\0A\00"
@.str4 = private global [16 x i8] c"DirectoryReader\00"
@.str5 = private global [22 x i8] c"assigining - %f to  \0A\00"
@.str6 = private global [9 x i8] c" - %s  \0A\00"
@.str7 = private global [2 x i8] c"a\00"
@.str8 = private global [22 x i8] c"assigining - %f to  \0A\00"
@.str9 = private global [9 x i8] c" - %s  \0A\00"
@.str10 = private global [2 x i8] c"z\00"
@.str11 = private global [18 x i8] c"return z = %f %d\0A\00"

define i64 @main.DirectoryReader.sum(i64 %x, %DirectoryReader* %this) {
entry:
	%0 = alloca i64
	store i64 %x, i64* %0
	%1 = getelementptr %DirectoryReader, %DirectoryReader* %this, i32 0, i32 0
	%2 = getelementptr %DirectoryReader, %DirectoryReader* %this, i32 0, i32 1
	%3 = getelementptr %Math, %Math* %2, i32 0, i32 0
	%4 = load double, double* %1
	%5 = load double, double* %3
	%6 = fadd double %4, %5
	%7 = alloca double
	store double %6, double* %7
	%8 = load double, double* %7
	%9 = fptosi double %8 to i64
	ret i64 %9
}

define i32 @main() {
entry:
	%0 = alloca %DirectoryReader
	%1 = alloca double
	store double 112.0, double* %1
	%2 = load double, double* %1
	%3 = getelementptr %DirectoryReader, %DirectoryReader* %0, i32 0, i32 0
	store double %2, double* %3
	%4 = alloca %Math
	%5 = alloca double
	store double 100.0, double* %5
	%6 = load double, double* %5
	%7 = getelementptr %Math, %Math* %4, i32 0, i32 0
	store double %6, double* %7
	%8 = getelementptr [15 x i8], [15 x i8]* @.str1, i64 0, i64 0
	%9 = getelementptr [5 x i8], [5 x i8]* @.str2, i64 0, i64 0
	%10 = call i32 (i8*, ...) @printf(i8* %8, i8* %9)
	%11 = load %Math, %Math* %4
	%12 = getelementptr %DirectoryReader, %DirectoryReader* %0, i32 0, i32 1
	store %Math %11, %Math* %12
	%13 = getelementptr [15 x i8], [15 x i8]* @.str3, i64 0, i64 0
	%14 = getelementptr [16 x i8], [16 x i8]* @.str4, i64 0, i64 0
	%15 = call i32 (i8*, ...) @printf(i8* %13, i8* %14)
	%16 = load %DirectoryReader, %DirectoryReader* %0
	%17 = getelementptr [22 x i8], [22 x i8]* @.str5, i64 0, i64 0
	%18 = call i32 (i8*, ...) @printf(i8* %17, %DirectoryReader %16)
	%19 = getelementptr [9 x i8], [9 x i8]* @.str6, i64 0, i64 0
	%20 = getelementptr [2 x i8], [2 x i8]* @.str7, i64 0, i64 0
	%21 = call i32 (i8*, ...) @printf(i8* %19, i8* %20)
	%22 = alloca double
	store double 10.0, double* %22
	%23 = load double, double* %22
	%24 = fptosi double %23 to i64
	%25 = call i64 @main.DirectoryReader.sum(i64 %24, %DirectoryReader* %0)
	%26 = alloca i64
	store i64 %25, i64* %26
	%27 = load i64, i64* %26
	%28 = getelementptr [22 x i8], [22 x i8]* @.str8, i64 0, i64 0
	%29 = call i32 (i8*, ...) @printf(i8* %28, i64 %27)
	%30 = getelementptr [9 x i8], [9 x i8]* @.str9, i64 0, i64 0
	%31 = getelementptr [2 x i8], [2 x i8]* @.str10, i64 0, i64 0
	%32 = call i32 (i8*, ...) @printf(i8* %30, i8* %31)
	%33 = alloca double
	store double 0.0, double* %33
	%34 = load double, double* %33
	%35 = fptosi double %34 to i32
	%36 = load i64, i64* %26
	%37 = getelementptr [18 x i8], [18 x i8]* @.str11, i64 0, i64 0
	%38 = call i32 (i8*, ...) @printf(i8* %37, i64 %36, i64 %36)
	ret i32 %35
}

declare i32 @printf(i8* %0, ...)
