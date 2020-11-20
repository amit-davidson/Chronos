package main

import (
	"bytes"
)

func main() {
	text := `Friends, Romans, countrymen, lend me your ears;
I come to bury Caesar, not to praise him.
The evil that men do lives after them;
The good is oft interred with their bones;
So let it be with Caesar. The noble Brutus
Hath told you Caesar was ambitious:
If it were so, it was a grievous fault,
And grievously hath Caesar answer'd it.
Here, under leave of Brutus and the rest -
For Brutus is an honourable man;
So are they all, all honourable men -
Come I to speak in Caesar's funeral.
He was my friend, faithful and just to me:
But Brutus says he was ambitious;
And Brutus is an honourable man.`

	data := bytes.NewBufferString(text)
	in := make(chan []byte)

	go func() {
		buf := make([]byte, 16)
		var n int
		var err error
		for ; err == nil; n, err = data.Read(buf) {
			in <- buf[:n]
		}
		close(in)
	}()
	res := ""
	for s := range in {
		res += string(s)
	}
	_ = res
}
