%DirectoryReader = type { double, %Math }
%Math = type { double }

define %DirectoryReader @main.DirectoryReader.sum(%Math %m, %DirectoryReader* %this) {
entry:
	%0 = alloca %Math
	store %Math %m, %Math* %0
	%1 = load %DirectoryReader, %DirectoryReader* %this
	ret %DirectoryReader %1
}

define {}* @main.DirectoryReader.math(%DirectoryReader* %this) {
entry:
	ret {}* null
}
