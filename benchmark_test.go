package iter

import (
	"strings"
	"testing"
	"unicode"
)

const benchString = `Go (often referred to as golang) is a programming language created at Google in 2009 by Robert Griesemer, Rob Pike, and Ken Thompson. It is a compiled, statically typed language in the tradition of Algol and C, with garbage collection, limited structural typing, memory safety features and CSP-style concurrent programming features added. The compiler and other language tools originally developed by Google are all free and open source.

Go (часто также Golang) — компилируемый многопоточный язык программирования, разработанный внутри компании Google. Первоначальная разработка Go началась в сентябре 2007 года, а его непосредственным проектированием занимались Роберт Гризмер, Роб Пайк и Кен Томпсон, занимавшиеся до этого проектом разработки операционной системы Inferno. Официально язык был представлен в ноябре 2009 года. На данный момент его поддержка осуществляется для операционных систем FreeBSD, OpenBSD, Linux, Mac OS X, Windows, начиная с версии 1.3 в язык Go включена экспериментальная поддержка DragonFly BSD, Plan 9 и Solaris, начиная с версии 1.4 — поддержка платформы Android.

غو (بالإنجليزية: GO) هي لغة برمجة مفتوحة المصدر من تطوير شركة جوجل. التصميم الأول للغة كان عام 2007 على يد روبرت غريسيمر و روب بايك و كِن ثومبسون. تم الإعلان رسمياً عن اللغة في نوفمبر 2009، مع تطبيقات صدرت لنظام التشغيل لينُكس و ماك. وقت صدورها، لم تعتبر جاهزة ليتم تبنيها في بيئات الإنتاج. في مايو 2010 صرح روب بايك علناً بأنه يتم استخدام اللغة لبعض الأمور المهمة في أنظمة جوجل.

Go是Google開發的一种静态强类型、編譯型、并发型，并具有垃圾回收功能的编程语言。為了方便搜索和識別，有时会将其稱為Golang。

羅伯特·格瑞史莫，羅勃·派克（Rob Pike）及肯·汤普逊於2007年9月开始设计Go語言，稍後Ian Lance Taylor、Russ Cox加入專案。Go語言是基於Inferno作業系統所開發的。Go語言於2009年11月正式宣布推出，成為開放原始碼專案，并在Linux及Mac OS X平台上进行了實現，后来追加了Windows系统下的实现。

目前Go语言每半年发布一个二级版本（即升级1.x到1.y）。

Go语言的语法接近C语言，但对于变量的声明有所不同。Go语言支持垃圾回收功能。Go语言的并行模型是以東尼·霍爾的交談循序程式（CSP）为基础，采取类似模型的其他语言包括Occam和Limbo，但它也具有Pi运算的特征，比如通道传输。在1.8版本中開放插件（Plugin）的支持，這意味著現在能從Go語言中動態載入部分函式。

与C++相比，Go語言並不包括如异常处理、继承、泛型、断言、虚函数等功能，但增加了 Slice 型、并发、管道、垃圾回收、接口（Interface）等特性的语言级支持。Google 目前仍正在討論是否應該支持泛型，其態度還是很開放的，但在該語言的常見問題列表中，對於断言的存在，則持負面態度，同時也為自己不提供型別繼承來辯護。

不同于Java，Go語言内嵌了关联数组（也称为哈希表（Hashes）或字典（Dictionaries）），就像字符串类型一样。`

var benchBytes = []byte(benchString)

func BenchmarkBytes(b *testing.B) {
	iter := New(benchBytes, spaceBytes)

	for i := 0; i < b.N; i++ {
		for iter.Next() {
			_ = iter.Bytes()
		}

		iter.Reset()
	}
}

func BenchmarkString(b *testing.B) {
	iter := NewString(benchString, space)

	for i := 0; i < b.N; i++ {
		for iter.Next() {
			_ = iter.String()
		}

		iter.Reset()
	}
}

func BenchmarkBytesEmitAll(b *testing.B) {
	iter := New(benchBytes, spaceBytes)
	iter.EmitAll()

	for i := 0; i < b.N; i++ {
		for iter.Next() {
			_ = iter.Bytes()
		}

		iter.Reset()
	}
}

func BenchmarkStringEmitAll(b *testing.B) {
	iter := NewString(benchString, space)
	iter.EmitAll()

	for i := 0; i < b.N; i++ {
		for iter.Next() {
			_ = iter.String()
		}

		iter.Reset()
	}
}

func BenchmarkBytesFunc(b *testing.B) {
	iter := NewFunc(benchBytes, unicode.IsSpace)

	for i := 0; i < b.N; i++ {
		for iter.Next() {
			_ = iter.Bytes()
		}

		iter.Reset()
	}
}

func BenchmarkStringFunc(b *testing.B) {
	iter := NewFuncString(benchString, unicode.IsSpace)

	for i := 0; i < b.N; i++ {
		for iter.Next() {
			_ = iter.String()
		}

		iter.Reset()
	}
}

func BenchmarkStdStringsMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = strings.FieldsFunc(benchString, unicode.IsSpace)
	}
}
