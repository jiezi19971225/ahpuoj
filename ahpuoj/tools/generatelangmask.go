package main

import (
	"fmt"
)

func main() {
	//C,Cpp,Pascal,Java,Ruby,Bash,Python,PHP,Perl,CSharp,Objc,FreeBasic,Schema,Clang,Clang++,Lua,Javascript,Go
	arr := []bool{/*C*/true,/*CPP*/true,/*pascal*/false,/*JAVA*/true,/*ruby*/false,/*bash*/false,/*python*/true,/*php*/false,/*perl*/false,/*csharp*/false,
	/*objc*/false,/*freebasic*/false,/*schema*/false,/*clang*/false,/*clang++*/false,/*lua*/false,/*js*/true,/*go*/true}
	res := 0
	for k,v := range(arr) {
		if v{
			res |= 1<<k
		}
	}
	fmt.Print(res)
}
