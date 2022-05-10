package benchmarks

import (
	"io"
	"testing"

	"github.com/rekby/genum/benchmarks/genumenum"
	"github.com/rekby/genum/benchmarks/stringerenum"
)

func BenchmarkSwitchIntEnum(b *testing.B) {
	values := []stringerenum.Enum{stringerenum.EnumA, stringerenum.EnumB, stringerenum.EnumC}
	sum := 0
	for i := 0; i < b.N; i++ {
		for index := range values {
			switch values[index] {
			case stringerenum.EnumA:
				sum += 1
			case stringerenum.EnumB:
				sum += 2
			case stringerenum.EnumC:
				sum += 3
			default:
				panic("unexpected")
			}
		}
	}
}

func BenchmarkSwitchGEnum(b *testing.B) {
	values := []genumenum.Variant{genumenum.A, genumenum.B, genumenum.C}
	sum := 0
	for i := 0; i < b.N; i++ {
		for index := range values {
			switch values[index] {
			case genumenum.A:
				sum += 1
			case genumenum.B:
				sum += 2
			case genumenum.C:
				sum += 3
			default:
				panic("unexpected")
			}
		}
	}
}

func BenchmarkStringIntWithStringer(b *testing.B) {
	v := stringerenum.EnumB
	values := make([]string, b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		values[i] = v.String()
	}

}

func BenchmarkStringGEnum(b *testing.B) {
	v := genumenum.B
	values := make([]string, b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		values[i] = v.String()
	}
}

func BenchmarkToStringGenumWithHolder(b *testing.B) {
	v := genumenum.B
	values := make([]string, b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		values[i] = genumenum.Holder.ValueToString(v)
	}
}

func BenchmarkFromIntInt(b *testing.B) {
	var v stringerenum.Enum
	for i := 0; i < b.N; i++ {
		v = stringerenum.Enum(1)
	}
	b.StopTimer()
	_, _ = io.Discard.Write([]byte(v.String()))
}

func BenchmarkFromIntGEnum(b *testing.B) {
	var v genumenum.Variant
	for i := 0; i < b.N; i++ {
		v, _ = genumenum.Holder.FromInt(1)
	}
	b.StopTimer()
	_, _ = io.Discard.Write([]byte(v.String()))
}

func BenchmarkFromIntGEnumUnsafe(b *testing.B) {
	var v genumenum.Variant
	for i := 0; i < b.N; i++ {
		v = genumenum.FromIntUnsafe(1)
	}
	b.StopTimer()
	_, _ = io.Discard.Write([]byte(v.String()))
}

func BenchmarkFromStringGEnum(b *testing.B) {
	var v genumenum.Variant
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v, _ = genumenum.Holder.FromString("A")
	}
	b.StopTimer()
	_, _ = io.Discard.Write([]byte(v.String()))

}
