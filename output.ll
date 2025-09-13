%DirectoryReader = type { double }
%Math = type { double }

@main.DirectoryReader.x = global double 200.0
@main.Math.x = global double 200.0
@.str1 = private global [14 x i8] c"final == %f \0A\00"

define i64 @main.DirectoryReader.sum(i64 %x, i64 %y, %DirectoryReader %this) {
entry:
	%0 = alloca i64
	store i64 %x, i64* %0
	%1 = alloca i64
	store i64 %y, i64* %1
	%2 = alloca %DirectoryReader
	%3 = getelementptr %DirectoryReader, %DirectoryReader* %2, i32 0, i32 0
	%4 = load double, double* %3
	%5 = alloca double
	store double %4, double* %5
	%6 = load double, double* %5
	%7 = fptosi double %6 to i64
	ret i64 %7
}

define i64 @main.DirectoryReader.mul(i64 %x, i64 %y, %DirectoryReader %this) {
entry:
	%0 = alloca i64
	store i64 %x, i64* %0
	%1 = alloca i64
	store i64 %y, i64* %1
	%2 = alloca %DirectoryReader
	%3 = load i64, i64* %0
	%4 = load i64, i64* %1
	%5 = sitofp i64 %3 to double
	%6 = sitofp i64 %4 to double
	%7 = fmul double %5, %6
	%8 = alloca double
	store double %7, double* %8
	%9 = load double, double* %8
	%10 = fptosi double %9 to i64
	ret i64 %10
}

define i64 @main.DirectoryReader.sub(i64 %x, i64 %y, %DirectoryReader %this) {
entry:
	%0 = alloca i64
	store i64 %x, i64* %0
	%1 = alloca i64
	store i64 %y, i64* %1
	%2 = alloca %DirectoryReader
	%3 = load i64, i64* %0
	%4 = load i64, i64* %1
	%5 = sitofp i64 %3 to double
	%6 = sitofp i64 %4 to double
	%7 = fsub double %5, %6
	%8 = alloca double
	store double %7, double* %8
	%9 = load double, double* %8
	%10 = fptosi double %9 to i64
	ret i64 %10
}

define i64 @main.Math.sum(i64 %x, i64 %y, %Math %this) {
entry:
	%0 = alloca i64
	store i64 %x, i64* %0
	%1 = alloca i64
	store i64 %y, i64* %1
	%2 = alloca %Math
	%3 = getelementptr %Math, %Math* %2, i32 0, i32 0
	%4 = load double, double* %3
	%5 = alloca double
	store double %4, double* %5
	%6 = load double, double* %5
	%7 = fptosi double %6 to i64
	ret i64 %7
}

define i64 @main.Math.mul(i64 %x, i64 %y, %Math %this) {
entry:
	%0 = alloca i64
	store i64 %x, i64* %0
	%1 = alloca i64
	store i64 %y, i64* %1
	%2 = alloca %Math
	%3 = load i64, i64* %0
	%4 = load i64, i64* %1
	%5 = sitofp i64 %3 to double
	%6 = sitofp i64 %4 to double
	%7 = fmul double %5, %6
	%8 = alloca double
	store double %7, double* %8
	%9 = load double, double* %8
	%10 = fptosi double %9 to i64
	ret i64 %10
}

define i64 @main.Math.sub(i64 %x, i64 %y, %Math %this) {
entry:
	%0 = alloca i64
	store i64 %x, i64* %0
	%1 = alloca i64
	store i64 %y, i64* %1
	%2 = alloca %Math
	%3 = load i64, i64* %0
	%4 = load i64, i64* %1
	%5 = sitofp i64 %3 to double
	%6 = sitofp i64 %4 to double
	%7 = fsub double %5, %6
	%8 = alloca double
	store double %7, double* %8
	%9 = load double, double* %8
	%10 = fptosi double %9 to i64
	ret i64 %10
}

define i32 @main() {
entry:
	%0 = alloca %DirectoryReader
	%1 = getelementptr %DirectoryReader, %DirectoryReader* %0, i32 0, i32 0
	store double 1.0, double* %1
	%2 = alloca double
	store double 10.0, double* %2
	%3 = load double, double* %2
	%4 = fptosi double %3 to i64
	%5 = alloca double
	store double 10.0, double* %5
	%6 = load double, double* %5
	%7 = fptosi double %6 to i64
	%8 = load %DirectoryReader, %DirectoryReader* %0
	%9 = call i64 @main.DirectoryReader.sum(i64 %4, i64 %7, %DirectoryReader %8)
	%10 = alloca i64
	store i64 %9, i64* %10
	%11 = alloca double
	store double 10.0, double* %11
	%12 = load double, double* %11
	%13 = fptosi double %12 to i64
	%14 = alloca double
	store double 20.0, double* %14
	%15 = load double, double* %14
	%16 = fptosi double %15 to i64
	%17 = load %DirectoryReader, %DirectoryReader* %0
	%18 = call i64 @main.DirectoryReader.mul(i64 %13, i64 %16, %DirectoryReader %17)
	%19 = alloca i64
	store i64 %18, i64* %19
	%20 = load i64, i64* %10
	%21 = load i64, i64* %19
	%22 = sitofp i64 %20 to double
	%23 = sitofp i64 %21 to double
	%24 = fadd double %22, %23
	%25 = alloca double
	store double %24, double* %25
	%26 = load double, double* %25
	%27 = getelementptr [14 x i8], [14 x i8]* @.str1, i64 0, i64 0
	%28 = call i32 (i8*, ...) @printf(i8* %27, double %26)
	ret i32 0
}

declare i32 @printf(i8* %0, ...)
