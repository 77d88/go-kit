package xbase64

import "testing"

func TestName(t *testing.T) {
	t.Log(Encode([]byte("hellohttp://abc!@$#@%*%()_)*!@#%world")))
	t.Log(EncodeURL([]byte("hellohttp://abc!@$#@%*%()_)*!@#%world")))
	t.Log(EncodeURL([]byte("hellohttp://abc!@$#@%*%()_)P{}LKLK^_^*③!@#%world")))
	t.Log(Encode([]byte("hellohttp://abc!@$#@%*%()_)P{}LKLK^_^*③!@#%world")))
	t.Log(RawURLEncode([]byte("hellohttp://abc!@$#@%*%()_)P{}LKLK^_^*③!@#%world")))
	t.Log(RawURLDecode("aGVsbG9odHRwOi8vYWJjIUAkI0AlKiUoKV8pUHt9TEtMS15fXirikaIhQCMld29ybGQ"))
}
