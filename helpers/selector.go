package helpers

/*
	This struct and it's methods / functions allows the user to answer 
	an output-input-request. Since there's almost no code, I'll not offer
	any documentation.

	You are still allowed to contact me via issues tho. :)
*/

import "fmt"

type Selector struct {
	Information []string
	Options     []string
}

func (s *Selector) SetInformation(information []string) {
	s.Information = s.Information[:0]
	s.Information = append(s.Information, information...)
}

func (s *Selector) AddOption(option string) {
	s.Options = append(s.Options, option)
}

func (s *Selector) Reset() {
	s.Information = s.Information[:0]
	s.Options = s.Options[:0]
}

func (s *Selector) RequestSelection() int {
	for _, info := range s.Information {
		fmt.Println(info)
	}

	fmt.Println("")

	for idx, option := range s.Options {
		str := fmt.Sprintf("> %d) %s", (idx + 1), option)
		fmt.Println(str)
	}

	fmt.Println("")

	fmt.Print("Choose your desired option: ")

	var i int
	fmt.Scan(&i)

	fmt.Println("")

	if i > len(s.Options) || (i < 1) {
		return -1
	}

	return i
}
