package testdata

type Something struct {
	ID int
	A  int
	B  string
	C  bool
	D  int64
}

func main() {
	//structinit:exhaustive
	var _ = Something{ // want `exhaustive struct literal .*structinit/testdata.Something not initialized with field D`
		ID: 1,
		A:  5,
		B:  "hello",
		C:  true,
	}

	var _ = Something{} // no diagnostic expected

	//structinit:exhaustive,omit=D
	var _ = Something{ // no diagnostic as D was omitted
		ID: 1,
		A:  5,
		B:  "hello",
		C:  true,
	}

	//structinit:exhaustive,omit=C,D
	var _ = Something{ // no diagnostic as C and D were omitted
		ID: 1,
		A:  5,
		B:  "hello",
	}

	//structinit:exhaustive,omit=Delta
	var _ = Something{ // want `omitted field Delta is not a field of .*structinit/testdata.Something`
		ID: 1,
		A:  5,
		B:  "hello",
		C:  true,
		D:  12,
	}

}
