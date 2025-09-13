%DirectoryReader = type { double }
%Main = type {}

@.DirectoryReader.x = global double 200.0
@.str1 = private global [16 x i8] c"param = %f, %s\0A\00"
@.str2 = private global [2 x i8] c"x\00"
@.str3 = private global [16 x i8] c"param = %f, %s\0A\00"
@.str4 = private global [2 x i8] c"y\00"
@.str5 = private global [16 x i8] c"param = %f, %s\0A\00"
@.str6 = private global [2 x i8] c"x\00"
@.str7 = private global [16 x i8] c"param = %f, %s\0A\00"
@.str8 = private global [2 x i8] c"y\00"
@.str9 = private global [16 x i8] c"param = %f, %s\0A\00"
@.str10 = private global [2 x i8] c"x\00"
@.str11 = private global [16 x i8] c"param = %f, %s\0A\00"
@.str12 = private global [2 x i8] c"y\00"
@.str13 = private global [14 x i8] c"final == %f \0A\00"

define i64 @.DirectoryReader.sum(i64 %x, i64 %y, %DirectoryReader %this) {
entry:
	%0 = alloca i64
	store i64 %x, i64* %0
	%1 = getelementptr [16 x i8], [16 x i8]* @.str1, i64 0, i64 0
	%2 = getelementptr [2 x i8], [2 x i8]* @.str2, i64 0, i64 0
	%3 = call i32 (i8*, ...) @printf(i8* %1, i64 %x, i8* %2)
	%4 = alloca i64
	store i64 %y, i64* %4
	%5 = getelementptr [16 x i8], [16 x i8]* @.str3, i64 0, i64 0
	%6 = getelementptr [2 x i8], [2 x i8]* @.str4, i64 0, i64 0
	%7 = call i32 (i8*, ...) @printf(i8* %5, i64 %y, i8* %6)
	%8 = alloca %DirectoryReader
	%9 = load i64, i64* %0
	%10 = load i64, i64* %4
	%11 = sitofp i64 %9 to double
	%12 = sitofp i64 %10 to double
	%13 = fadd double %11, %12
	%14 = alloca double
	store double %13, double* %14
	%15 = load double, double* %14
	%16 = fptosi double %15 to i64
	ret i64 %16
}

define i64 @.DirectoryReader.mul(i64 %x, i64 %y, %DirectoryReader %this) {
entry:
	%0 = alloca i64
	store i64 %x, i64* %0
	%1 = getelementptr [16 x i8], [16 x i8]* @.str5, i64 0, i64 0
	%2 = getelementptr [2 x i8], [2 x i8]* @.str6, i64 0, i64 0
	%3 = call i32 (i8*, ...) @printf(i8* %1, i64 %x, i8* %2)
	%4 = alloca i64
	store i64 %y, i64* %4
	%5 = getelementptr [16 x i8], [16 x i8]* @.str7, i64 0, i64 0
	%6 = getelementptr [2 x i8], [2 x i8]* @.str8, i64 0, i64 0
	%7 = call i32 (i8*, ...) @printf(i8* %5, i64 %y, i8* %6)
	%8 = alloca %DirectoryReader
	%9 = load i64, i64* %0
	%10 = load i64, i64* %4
	%11 = sitofp i64 %9 to double
	%12 = sitofp i64 %10 to double
	%13 = fmul double %11, %12
	%14 = alloca double
	store double %13, double* %14
	%15 = load double, double* %14
	%16 = fptosi double %15 to i64
	ret i64 %16
}

define i64 @.DirectoryReader.sub(i64 %x, i64 %y, %DirectoryReader %this) {
entry:
	%0 = alloca i64
	store i64 %x, i64* %0
	%1 = getelementptr [16 x i8], [16 x i8]* @.str9, i64 0, i64 0
	%2 = getelementptr [2 x i8], [2 x i8]* @.str10, i64 0, i64 0
	%3 = call i32 (i8*, ...) @printf(i8* %1, i64 %x, i8* %2)
	%4 = alloca i64
	store i64 %y, i64* %4
	%5 = getelementptr [16 x i8], [16 x i8]* @.str11, i64 0, i64 0
	%6 = getelementptr [2 x i8], [2 x i8]* @.str12, i64 0, i64 0
	%7 = call i32 (i8*, ...) @printf(i8* %5, i64 %y, i8* %6)
	%8 = alloca %DirectoryReader
	%9 = load i64, i64* %0
	%10 = load i64, i64* %4
	%11 = sitofp i64 %9 to double
	%12 = sitofp i64 %10 to double
	%13 = fsub double %11, %12
	%14 = alloca double
	store double %13, double* %14
	%15 = load double, double* %14
	%16 = fptosi double %15 to i64
	ret i64 %16
}

define i64 @.Main.creat(%Main %this) {
entry:
	%0 = alloca %Main
	%1 = alloca double
	store double 0.0, double* %1
	%2 = load double, double* %1
	%3 = fptosi double %2 to i64
	ret i64 %3
}

declare i32 @printf(i8* %0, ...)

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
	%9 = call i64 @.DirectoryReader.sum(i64 %4, i64 %7, %DirectoryReader %8)
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
	%18 = call i64 @.DirectoryReader.mul(i64 %13, i64 %16, %DirectoryReader %17)
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
	%27 = getelementptr [14 x i8], [14 x i8]* @.str13, i64 0, i64 0
	%28 = call i32 (i8*, ...) @printf(i8* %27, double %26)
	ret i32 0
}
