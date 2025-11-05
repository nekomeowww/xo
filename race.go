package xo

func Race2[T0 any, T1 any](src1 chan T0, src2 chan T1) (T0, T1) {
	var empty0 T0
	var empty1 T1

	select {
	case v := <-src1:
		return v, empty1
	case v := <-src2:
		return empty0, v
	}
}



func Race3[T0 any, T1 any, T2 any](src1 chan T0, src2 chan T1, src3 chan T2) (T0, T1, T2) {
	var empty0 T0
	var empty1 T1
	var empty2 T2

	select {
	case v := <-src1:
		return v, empty1, empty2
	case v := <-src2:
		return empty0, v, empty2
	case v := <-src3:
		return empty0, empty1, v
	}
}



func Race4[T0 any, T1 any, T2 any, T3 any](src1 chan T0, src2 chan T1, src3 chan T2, src4 chan T3) (T0, T1, T2, T3) {
var empty0 T0
var empty1 T1
var empty2 T2
var empty3 T3

select {
case v := <-src1:
return v, empty1, empty2, empty3
case v := <-src2:
return empty0, v, empty2, empty3
case v := <-src3:
return empty0, empty1, v, empty3
case v := <-src4:
return empty0, empty1, empty2, v
}
}



func Race5[T0 any, T1 any, T2 any, T3 any, T4 any](src1 chan T0, src2 chan T1, src3 chan T2, src4 chan T3, src5 chan T4) (T0, T1, T2, T3, T4) {
	var empty0 T0
	var empty1 T1
	var empty2 T2
	var empty3 T3
	var empty4 T4

	select {
	case v := <-src1:
		return v, empty1, empty2, empty3, empty4
	case v := <-src2:
		return empty0, v, empty2, empty3, empty4
	case v := <-src3:
		return empty0, empty1, v, empty3, empty4
	case v := <-src4:
		return empty0, empty1, empty2, v, empty4
	case v := <-src5:
		return empty0, empty1, empty2, empty3, v
	}
}



func Race6[T0 any, T1 any, T2 any, T3 any, T4 any, T5 any](src1 chan T0, src2 chan T1, src3 chan T2, src4 chan T3, src5 chan T4, src6 chan T5) (T0, T1, T2, T3, T4, T5) {
	var empty0 T0
	var empty1 T1
	var empty2 T2
	var empty3 T3
	var empty4 T4
	var empty5 T5

	select {
	case v := <-src1:
		return v, empty1, empty2, empty3, empty4, empty5
	case v := <-src2:
		return empty0, v, empty2, empty3, empty4, empty5
	case v := <-src3:
		return empty0, empty1, v, empty3, empty4, empty5
	case v := <-src4:
		return empty0, empty1, empty2, v, empty4, empty5
	case v := <-src5:
		return empty0, empty1, empty2, empty3, v, empty5
	case v := <-src6:
		return empty0, empty1, empty2, empty3, empty4, v
	}
}

