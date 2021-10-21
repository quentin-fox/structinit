package test

import "external"

type Something struct {
	ID int
	A  int
	B  string
	C  bool
	D  int64
}

func main() {
	//structinit:exhaustive
	var _ = Something{ // want `exhaustive struct literal test.Something not initialized with field D`
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
	var _ = Something{ // want `omitted field Delta is not a field of test.Something`
		ID: 1,
		A:  5,
		B:  "hello",
		C:  true,
		D:  12,
	}

	// structinit shouldn't complain about not initializing structs with private fields that this package doesn't have access to

	var _ = external.Something{} // no diagnostic required


	//structinit:exhaustive
	var _ = external.Something{} // want `exhaustive struct literal external.Something not initialized with field ID`


	//structinit:exhaustive
	var _ = external.Something{ // no diagnostic required
		ID: 15,
	}

	//structinit:exhaustive,omit=ID
	var _ = external.Something{} // no diagnost required
}
